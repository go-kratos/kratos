package errors

import (
	"fmt"
	"go-common/library/ecode"
)

// DependError for show like "model(%s) code(%s) error"
type DependError struct {
	errInfo string
	Err     error
}

// New new a video error.
func New(errInfo string, m error) error {
	return &DependError{errInfo, m}
}

// Ecode return ecode.
func (e *DependError) Ecode() error {
	return e.Err
}

// Error implement error interface.
func (e *DependError) Error() string {
	m := ecode.String(e.Err.Error()).Message()
	return fmt.Sprintf(m, e.errInfo)
}
