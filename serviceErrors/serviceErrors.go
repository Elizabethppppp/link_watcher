package serviceErrors

import "errors"

var (
	ErrBadRequest      = errors.New("bad request")
	ErrInvalidId       = errors.New("invalid id")
	ErrInvalidJSON     = errors.New("invalid json")
	ErrInvalidURL      = errors.New("invalid url")
	ErrURLTooLarge     = errors.New("url too large")
	ErrEmptyURL        = errors.New("url is empty")
	ErrInvalidInterval = errors.New("invalid interval")
	ErrNotFound        = errors.New("not found")
	ErrConflict        = errors.New("conflict")
	ErrInternal        = errors.New("internal error")
)
