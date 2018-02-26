package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInfo(t *testing.T) {
	labels := map[string]string{
		"foo": "bar",
	}

	info := NewResolver(labels).GetInfo(5000, "test-version")

	assert.NotEmpty(t, info.BootID, "should resolve BootID")
	assert.NotEmpty(t, info.MachineID, "should resolve MachineID")
	assert.NotEmpty(t, info.SystemUUID, "should resolve SystemUUID")
	assert.Equal(t, labels, info.Labels, "should have given node labels")
	assert.Equal(t, 5000, info.GrpcPort, "should have given node grpc port")
}
