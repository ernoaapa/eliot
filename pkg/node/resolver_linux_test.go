package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInfo(t *testing.T) {
	labels := map[string]string{
		"foo": "bar",
	}

	info := NewResolver(5000, "test-version", labels).GetInfo()

	assert.NotEmpty(t, info.BootID, "should resolve BootID")
	assert.NotEmpty(t, info.MachineID, "should resolve MachineID")
	assert.NotEmpty(t, info.SystemUUID, "should resolve SystemUUID")
	assert.Equal(t, labels, info.Labels, "should have given node labels")
	assert.Equal(t, 5000, info.GrpcPort, "should have given node grpc port")
	assert.True(t, len(info.Filesystems) > 0, "should have at least one disk")
}

func TestResolveFilesystems(t *testing.T) {
	// Warning: we're assuming that we run in environment where is filesystem info is available
	assert.True(t, len(resolveFilesystems()) > 0)
}

func TestResolveUptime(t *testing.T) {
	// Warning: we're assuming that we run in environment where is uptime info is available
	assert.True(t, resolveUptime() > 0)
}
