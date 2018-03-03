package runtime

import (
	"fmt"

	"github.com/pkg/errors"
)

// Definitions of common error types used throughout runtime implementation.
// All errors returned by the interface will map into one of these errors classes.
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrNotSupported  = errors.New("not supported")
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
