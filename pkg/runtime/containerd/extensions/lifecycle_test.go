package extensions

import (
	"github.com/containerd/containerd/containers"
	"github.com/containerd/typeurl"
	"testing"

	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
)

func TestGetLifecycleExtension(t *testing.T) {
	lifecycle := &ContainerLifecycle{
		StartCount:    666,
		RestartPolicy: OnFailure,
	}
	any, _ := typeurl.MarshalAny(lifecycle)
	extensions := make(map[string]types.Any)
	extensions[lifecycleExtensionName] = *any

	result, err := GetLifecycleExtension(containers.Container{Extensions: extensions})
	assert.NoError(t, err)
	assert.Equal(t, lifecycle.StartCount, result.StartCount)
	assert.Equal(t, lifecycle.RestartPolicy, result.RestartPolicy)
}

func TestGetLifecycleExtensionReturnNotFoundErr(t *testing.T) {
	_, err := GetLifecycleExtension(containers.Container{})
	assert.True(t, IsNotFound(err))
}
