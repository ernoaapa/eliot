package runtime

import (
	"testing"

	imagespecs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
)

func TestPlatformExist(t *testing.T) {
	assert.True(t, platformExist("linux/arm64", []imagespecs.Platform{
		{OS: "linux", Architecture: "arm64"},
	}))

	assert.False(t, platformExist("linux/arm64", []imagespecs.Platform{
		{OS: "linux", Architecture: "amd64"},
	}))
}
