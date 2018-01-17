package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	sync, err := Parse("~/local/dir:/data")
	assert.NoError(t, err)

	assert.Equal(t, "~/local/dir", sync.Source)
	assert.Equal(t, "/data", sync.Destination)
}
