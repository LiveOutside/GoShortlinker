package activationcodes

import "errors"

var (
	ErrInternal     = errors.New("internal server error")
	ErrCodeNotFound = errors.New("activation code not found")
	ErrCodeExpired  = errors.New("activation code expired")
	ErrInvalidCode  = errors.New("invalid code provided")
)
