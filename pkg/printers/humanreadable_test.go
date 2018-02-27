package printers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatPercent(t *testing.T) {
	assert.Equal(t, "80%", formatPercent(100*1024, 20*1024, 20*1024))
}
