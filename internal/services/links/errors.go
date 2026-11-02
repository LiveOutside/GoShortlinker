package links

import "errors"

var (
	ErrSaveAndShortenLink = errors.New("failed to shorten link")
	ErrInvalidParams      = errors.New("invalid link params")
	ErrLinkNotFound       = errors.New("link not found")
	ErrInternal           = errors.New("internal error")
	ErrLinkUnavailable    = errors.New("link is not accessable")
)
