package humanreadable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDotsSpinner(t *testing.T) {
	spinner := NewDots()
	spinner2 := NewDots()
	assert.Equal(t, `⠙`, spinner.Rotate())
	assert.Equal(t, `⠹`, spinner.Rotate())

	assert.Equal(t, `⠙`, spinner2.Rotate())
}
