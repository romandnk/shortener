package urlservice

import "errors"

var (
	ErrInternalError = errors.New("internal error")
)

var (
	ErrInvalidOriginalURL = errors.New("invalid url format")
	ErrEmptyOriginalURL   = errors.New("url cannot be empty")
	ErrOriginalURLTooLong = errors.New("max url length is 2048")
)
