package handlers

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"os"
	"strconv"
	"sync"

	"goshortlinker/internal/middleware"
	authservice "goshortlinker/internal/services/auth"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"golang.org/x/oauth2"
)

type socialProvider struct {
	cfg          *oauth2.Config
	exchangeOpts func(c fiber.Ctx) []oauth2.AuthCodeOption
	identify     func(ctx context.Context, tok *oauth2.Token) (string, error)
}

type SocialHandler struct {
	auth        *authservice.Service
	authHandler *AuthenticationHandler
	providers   map[string]*socialProvider
	secure      bool
}

func NewSocialHandler(s *authservice.Service, ah *AuthenticationHandler) *SocialHandler {
	base := os.Getenv("APP_BASE_URL")

	vk := &socialProvider{
		cfg: &oauth2.Config{
			ClientID:     os.Getenv("VK_CLIENT_ID"),
			ClientSecret: os.Getenv("VK_CLIENT_SECRET"),
			RedirectURL:  base + "/auth/vk/callback",
			Scopes:       []string{"vkid.personal_info"},
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://id.vk.com/authorize",
				TokenURL:  "https://id.vk.com/oauth2/auth",
				AuthStyle: oauth2.AuthStyleInParams,
			},
		},
		exchangeOpts: func(c fiber.Ctx) []oauth2.AuthCodeOption {
			return []oauth2.AuthCodeOption{
				oauth2.SetAuthURLParam("device_id", c.Query("device_id")),
				oauth2.SetAuthURLParam("state", c.Query("state")),
			}
		},
		identify: func(_ context.Context, tok *oauth2.Token) (string, error) {
			id := idToString(tok.Extra("user_id"))
			if id == "" {
				return "", errors.New("vk: no user_id in token response")
			}
			return id, nil
		},
	}

	tgClientID := os.Getenv("TELEGRAM_CLIENT_ID")
	var (
		tgMu sync.Mutex
		tgV  *oidc.IDTokenVerifier
	)
	tgVerifier := func() (*oidc.IDTokenVerifier, error) {
		tgMu.Lock()
		defer tgMu.Unlock()
		if tgV != nil {
			return tgV, nil
		}
		p, err := oidc.NewProvider(context.Background(), "https://oauth.telegram.org")
		if err != nil {
			return nil, err
		}
		tgV = p.Verifier(&oidc.Config{ClientID: tgClientID})
		return tgV, nil
	}
	tg := &socialProvider{
		cfg: &oauth2.Config{
			ClientID:     tgClientID,
			ClientSecret: os.Getenv("TELEGRAM_CLIENT_SECRET"),
			RedirectURL:  base + "/auth/telegram/callback",
			Scopes:       []string{"openid"},
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://oauth.telegram.org/auth",
				TokenURL:  "https://oauth.telegram.org/token",
				AuthStyle: oauth2.AuthStyleInHeader,
			},
		},
		exchangeOpts: func(fiber.Ctx) []oauth2.AuthCodeOption {
			return []oauth2.AuthCodeOption{oauth2.SetAuthURLParam("client_id", tgClientID)}
		},
		identify: func(ctx context.Context, tok *oauth2.Token) (string, error) {
			raw, _ := tok.Extra("id_token").(string)
			if raw == "" {
				return "", errors.New("telegram: no id_token")
			}
			v, err := tgVerifier()
			if err != nil {
				return "", err
			}
			idt, err := v.Verify(ctx, raw)
			if err != nil {
				return "", err
			}
			return idt.Subject, nil
		},
	}

	return &SocialHandler{
		auth:        s,
		authHandler: ah,
		providers:   map[string]*socialProvider{"vk": vk, "telegram": tg},
		secure:      os.Getenv("COOKIE_SECURE") != "false",
	}
}

func (h *SocialHandler) StartLogin(c fiber.Ctx) error { return h.start(c, "login") }

func (h *SocialHandler) StartLink(c fiber.Ctx) error { return h.start(c, "link") }

func (h *SocialHandler) start(c fiber.Ctx, intent string) error {
	p, ok := h.providers[c.Params("provider")]
	if !ok {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if intent == "link" {
		if _, ok := middleware.UserID(c); !ok {
			return c.Redirect().To("/auth")
		}
	}

	state := randomString(32)
	verifier := oauth2.GenerateVerifier()
	h.setTmp(c, "oauth_state", state)
	h.setTmp(c, "oauth_verifier", verifier)
	h.setTmp(c, "oauth_intent", intent)

	return c.Redirect().To(p.cfg.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier)))
}

func (h *SocialHandler) Callback(c fiber.Ctx) error {
	name := c.Params("provider")
	p, ok := h.providers[name]
	if !ok {
		return c.SendStatus(fiber.StatusNotFound)
	}

	state := c.Cookies("oauth_state")
	verifier := c.Cookies("oauth_verifier")
	intent := c.Cookies("oauth_intent")
	for _, n := range []string{"oauth_state", "oauth_verifier", "oauth_intent"} {
		h.clearTmp(c, n)
	}

	back := "/auth"
	if intent == "link" {
		back = "/profile"
	}
	if (intent != "login" && intent != "link") || state == "" || verifier == "" ||
		subtle.ConstantTimeCompare([]byte(state), []byte(c.Query("state"))) != 1 ||
		c.Query("error") != "" || c.Query("code") == "" {
		return c.Redirect().To(back + "?error=oauth")
	}

	ctx := c.Context()
	opts := []oauth2.AuthCodeOption{oauth2.VerifierOption(verifier)}
	if p.exchangeOpts != nil {
		opts = append(opts, p.exchangeOpts(c)...)
	}
	tok, err := p.cfg.Exchange(ctx, c.Query("code"), opts...)
	if err != nil {
		log.Errorf("%s exchange: %v", name, err)
		return c.Redirect().To(back + "?error=oauth")
	}
	providerID, err := p.identify(ctx, tok)
	if err != nil {
		log.Errorf("%s identify: %v", name, err)
		return c.Redirect().To(back + "?error=oauth")
	}

	if intent == "link" {
		uid, ok := middleware.UserID(c) // сессия пользователя из JWT-cookie
		if !ok {
			return c.Redirect().To("/auth")
		}
		switch err := h.auth.LinkSocial(ctx, uid, name, providerID); {
		case errors.Is(err, authservice.ErrSocialTaken):
			return c.Redirect().To("/profile?error=social_taken")
		case err != nil:
			log.Errorf("link %s: %v", name, err)
			return c.Redirect().To("/profile?error=server")
		}
		return c.Redirect().To("/profile")
	}

	res, err := h.auth.LoginSocial(ctx, name, providerID)
	switch {
	case errors.Is(err, authservice.ErrSocialNotLinked):
		return c.Redirect().To("/auth?error=social_not_linked")
	case errors.Is(err, authservice.ErrNotActivated):
		return c.Redirect().To("/activate?user_id=" + strconv.Itoa(int(res.UserID)))
	case err != nil:
		log.Errorf("login %s: %v", name, err)
		return c.Redirect().To("/auth?error=server")
	}
	h.authHandler.setAuthCookie(c, res.Token, res.ExpiresAt)
	return c.Redirect().To("/")
}

func (h *SocialHandler) Unlink(c fiber.Ctx) error {
	uid, ok := middleware.UserID(c)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	name := c.Params("provider")
	if _, ok := h.providers[name]; !ok {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if err := h.auth.UnlinkSocial(c.Context(), uid, name); err != nil {
		log.Errorf("unlink %s: %v", name, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "внутренняя ошибка сервера"})
	}
	return c.JSON(fiber.Map{"ok": true})
}

func (h *SocialHandler) setTmp(c fiber.Ctx, name, value string) {
	c.Cookie(&fiber.Cookie{
		Name: name, Value: value, Path: "/", MaxAge: 600,
		HTTPOnly: true, Secure: h.secure, SameSite: fiber.CookieSameSiteLaxMode,
	})
}

func (h *SocialHandler) clearTmp(c fiber.Ctx, name string) {
	c.Cookie(&fiber.Cookie{
		Name: name, Value: "", Path: "/", MaxAge: -1,
		HTTPOnly: true, Secure: h.secure, SameSite: fiber.CookieSameSiteLaxMode,
	})
}

func randomString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func idToString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatInt(int64(t), 10)
	}
	return ""
}
