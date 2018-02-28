package printers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatPercent(t *testing.T) {
	assert.Equal(t, "80%", formatPercent(100*1024, 20*1024, 20*1024))
}

func TestFormatDuration(t *testing.T) {
	assert.Equal(t, "292 years 24 weeks 3 days 23 hours 47 minutes 16 seconds 854 milliseconds", formatUptime(9223372036854775807), "should format maximum int64 (half of uint64)")
	assert.Equal(t, "-", formatUptime(18446744073709551615), "Should not break if goes above int64 (e.g. if maximum uint64)")
}
