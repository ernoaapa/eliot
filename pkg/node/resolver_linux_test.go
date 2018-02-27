package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveFilesystems(t *testing.T) {
	// Warning: we're assuming that we run in environment where is filesystem info is available
	assert.True(t, len(resolveFilesystems()) > 0)
}
