package device

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInfo(t *testing.T) {
	info := GetInfo()
	log.Fatal(info.MachineID)
	assert.NotEmpty(t, info.BootID, "should resolve BootID")
	assert.NotEmpty(t, info.MachineID, "should resolve MachineID")
	assert.NotEmpty(t, info.SystemUUID, "should resolve SystemUUID")
}
