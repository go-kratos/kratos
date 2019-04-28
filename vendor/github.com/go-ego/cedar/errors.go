package cedar

import (
	"errors"
)

var (
	// ErrInvalidDataType invalid data type error
	ErrInvalidDataType = errors.New("cedar: invalid datatype")
	// ErrInvalidValue invalid value error
	ErrInvalidValue = errors.New("cedar: invalid value")
	// ErrInvalidKey invalid key error
	ErrInvalidKey = errors.New("cedar: invalid key")
	// ErrNoPath no path error
	ErrNoPath = errors.New("cedar: no path")
	// ErrNoValue no value error
	ErrNoValue = errors.New("cedar: no value")
)
