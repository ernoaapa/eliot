package device

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInfo(t *testing.T) {
	labels := map[string]string{
		"foo": "bar",
	}

	info := NewResolver(labels).GetInfo()

	assert.NotEmpty(t, info.BootID, "should resolve BootID")
	assert.NotEmpty(t, info.MachineID, "should resolve MachineID")
	assert.NotEmpty(t, info.SystemUUID, "should resolve SystemUUID")
	assert.Equal(t, labels, info.Labels, "should have given device labels")
}
