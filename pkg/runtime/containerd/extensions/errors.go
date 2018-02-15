package extensions

import (
	"fmt"

	"github.com/pkg/errors"
)

// Definitions of common error types from extensions
var (
	ErrNotFound = errors.New("not found")
)

// IsNotFound returns true if the error is due to a missing resource
func IsNotFound(err error) bool {
	return errors.Cause(err) == ErrNotFound
}

// ErrWithMessagef updates error message with formated message
// I.e. errors.WithMessage(err, fmt.Sprintf(...
// Hopefully we can change to errors.WithMessagef some day: https://github.com/pkg/errors/pull/118
func ErrWithMessagef(err error, format string, args ...interface{}) error {
	return errors.WithMessage(err, fmt.Sprintf(format, args...))
}
