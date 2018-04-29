package runtime

import (
	"testing"

	"github.com/containerd/containerd/platforms"
	imagespecs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
)

func TestPlatformExist(t *testing.T) {
	platform, err := platforms.Parse("linux/arm64")
	assert.NoError(t, err)

	assert.True(t, platformExist(platform, []imagespecs.Platform{
		{OS: "linux", Architecture: "arm64"},
	}))

	assert.False(t, platformExist(platform, []imagespecs.Platform{
		{OS: "linux", Architecture: "amd64"},
	}))
}
