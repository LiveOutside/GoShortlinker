package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrNotActivated       = errors.New("account is not activated")
	ErrInvalidToken       = errors.New("invalid token")
)
