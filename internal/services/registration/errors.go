package users

import "errors"

var (
	ErrPasswordMismatch      = errors.New("passwords do not match")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInternal              = errors.New("internal error")
)
