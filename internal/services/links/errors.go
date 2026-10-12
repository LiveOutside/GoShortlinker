package links

import "errors"

var (
	ErrSaveAndShortenLink = errors.New("failed to shorten link")
	ErrInvalidParams      = errors.New("invalid link params")
)
