package handlers

import (
	"errors"

	dtousers "goshortlinker/internal/data/dto/users"
	"goshortlinker/internal/middleware"
	authservice "goshortlinker/internal/services/auth"
	registrationservice "goshortlinker/internal/services/registration"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

var profileErrors = map[string]string{
	"social_taken": "Эта соцсеть уже привязана к другому пользователю",
	"oauth":        "Не удалось привязать соцсеть, попробуйте ещё раз",
	"server":       "Внутренняя ошибка сервера",
}

type ProfileHandler struct {
	profile     *registrationservice.Service
	auth        *authservice.Service
	authHandler *AuthenticationHandler
}

func NewProfileHandler(p *registrationservice.Service, a *authservice.Service, ah *AuthenticationHandler) *ProfileHandler {
	return &ProfileHandler{profile: p, auth: a, authHandler: ah}
}

func (h *ProfileHandler) Get(c fiber.Ctx) error {
	userID, ok := middleware.UserID(c)
	if !ok {
		return c.Redirect().To("/auth")
	}

	profile, err := h.profile.GetProfile(c.Context(), userID)
	if err != nil {
		if errors.Is(err, registrationservice.ErrUserNotFound) {
			h.authHandler.clearAuthCookie(c)
			return c.Redirect().To("/auth")
		}
		log.Errorf("get profile: %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString("внутренняя ошибка сервера")
	}

	linked, err := h.auth.LinkedSocials(c.Context(), userID)
	if err != nil {
		log.Errorf("linked socials: %v", err)
		linked = map[string]bool{}
	}

	return c.Render("profile", fiber.Map{
		"Profile": profile,
		"Socials": []fiber.Map{
			{"Key": "vk", "Title": "VK", "Linked": linked["vk"]},
			{"Key": "telegram", "Title": "Telegram", "Linked": linked["telegram"]},
		},
		"SocialError": profileErrors[c.Query("error")],
	})
}

func (h *ProfileHandler) Update(c fiber.Ctx) error {
	userID, ok := middleware.UserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":       "требуется авторизация",
			"redirect_to": "/auth",
		})
	}

	var payload dtousers.UpdateProfileRequest
	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "некорректный запрос"})
	}
	payload.Normalize()
	if err := validate.Struct(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": profileValidationMessage(err)})
	}

	updated, err := h.profile.UpdateProfile(c.Context(), userID, payload)
	switch {
	case errors.Is(err, registrationservice.ErrUsernameAlreadyExists):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "это имя пользователя уже занято"})
	case errors.Is(err, registrationservice.ErrEmailAlreadyExists):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "этот email уже используется"})
	case errors.Is(err, registrationservice.ErrUserNotFound):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "пользователь не найден", "redirect_to": "/auth"})
	case err != nil:
		log.Errorf("update profile: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "внутренняя ошибка сервера"})
	}

	if res, err := h.auth.Reissue(userID, updated.Username); err != nil {
		log.Errorf("reissue token: %v", err)
	} else {
		h.authHandler.setAuthCookie(c, res.Token, res.ExpiresAt)
	}

	return c.JSON(updated)
}

func profileValidationMessage(err error) string {
	var ves validator.ValidationErrors
	if errors.As(err, &ves) && len(ves) > 0 {
		switch ves[0].Field() {
		case "Username":
			return "Имя пользователя: 3–32 символа, только латинские буквы и цифры"
		case "Email":
			return "Введите корректный email (не длиннее 255 символов)"
		case "Title":
			return "Должность — не более 100 символов"
		case "Bio":
			return "О себе — не более 500 символов"
		}
	}
	return "некорректные данные"
}
