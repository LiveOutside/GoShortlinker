package users

import "errors"

var ErrPasswordMismatch = errors.New("passwords do not match")
var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrUsernameAlreadyExists = errors.New("username already exists")
var ErrInternal = errors.New("internal error")
