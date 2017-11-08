package runtime

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestIsNotFound(t *testing.T) {
	assert.True(t, IsNotFound(ErrNotFound))
	assert.True(t, IsNotFound(errors.Wrapf(ErrNotFound, "Foo bar not found")), "should support custom message")
	assert.True(t, IsNotFound(ErrWithMessagef(ErrNotFound, "Foo bar not found")))
	assert.False(t, IsNotFound(ErrWithMessagef(ErrAlreadyExists, "Foo bar not found")), "should not pass if not ErrNotFound")
}
