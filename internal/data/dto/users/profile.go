package users

import "strings"

type UpdateProfileRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32,alphanum"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Title    string `json:"title" validate:"max=100"`
	Bio      string `json:"bio" validate:"max=500"`
}

func (r *UpdateProfileRequest) Normalize() {
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.Title = strings.TrimSpace(r.Title)
	r.Bio = strings.TrimSpace(r.Bio)
}

type ProfileResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Title    string `json:"title"`
	Bio      string `json:"bio"`
}
