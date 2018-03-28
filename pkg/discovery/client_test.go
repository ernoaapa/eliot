package discovery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClientNodes(t *testing.T) {
	server := NewServer("testing", 1234, "v1.0")
	go server.Serve()
	defer server.Stop()

	nodes, err := Nodes(1 * time.Second)
	assert.NoError(t, err)
	assert.Len(t, nodes, 1)
	assert.Equal(t, int64(1234), nodes[0].GrpcPort)
	assert.Equal(t, "v1.0", nodes[0].Version)
}
