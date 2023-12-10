package storageerrors

import "errors"

var ErrInvalidDB = errors.New("invalid db type")

var (
	ErrOriginalURLExists = errors.New("original url already exists")
	ErrURLAliasExists    = errors.New("url alias already exists")
	ErrShortNotFound     = errors.New("short urs is not found")
)
