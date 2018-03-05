package printers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatPercent(t *testing.T) {
	assert.Equal(t, "80%", formatPercent(100*1024, 20*1024, 20*1024))
}

func TestFormatDuration(t *testing.T) {
	assert.Equal(t, "2 minutes 40 seconds", formatUptime(160), "Should format 160 seconds")
	assert.Equal(t, "292 years 24 weeks 3 days 23 hours 47 minutes 16 seconds", formatUptime(9223372036), "should format large value (maximum Nanosecond duration in seconds)")
	assert.Equal(t, "18446744073709551615 seconds", formatUptime(18446744073709551615), "Should not break if goes above int64 (e.g. if maximum uint64)")
}
