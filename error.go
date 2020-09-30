package configurator

import (
	"errors"
)

var (
	ErrInvalidConfig    = errors.New("config must be a struct pointer")
	ErrInvalidTagFormat = errors.New("invalid tag format")
	ErrEmptyValue       = errors.New("empty value")
	ErrEmptyKey         = errors.New("empty key")
	ErrConflictKey      = errors.New("conflict key")
	ErrUnsupported      = errors.New("unsupported")
)
