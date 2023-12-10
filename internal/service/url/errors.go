package urlservice

import "errors"

var (
	ErrInternalError = errors.New("internal error")
)

var (
	ErrInvalidOriginalURL = errors.New("invalid url format")
	ErrEmptyOriginalURL   = errors.New("url cannot be empty")
	ErrOriginalURLTooLong = errors.New("max url length is 2048")

	ErrEmptyURLAlias       = errors.New("empty url unique id")
	ErrInvalidAliasFormat  = errors.New("unique id must be 10 symbols length")
	ErrOriginalURLNotFound = errors.New("original url is not found")
)
