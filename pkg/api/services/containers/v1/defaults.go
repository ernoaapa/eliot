package containers

import (
	"io/ioutil"
	"path/filepath"
)

// Defaults set default values to container definitions
func Defaults(containers []*Container) (result []*Container) {
	for _, container := range containers {
		result = append(result, Default(container))
	}
	return result
}

// Default set default values to Container model
func Default(container *Container) *Container {
	if container.Io == nil {
		container.Io = &IOSet{}
	}

	dir, _ := ioutil.TempDir("/run/containerd/fifo", "")
	if container.Io.In == "" {
		container.Io.In = filepath.Join(dir, "-stdin")
	}
	if container.Io.Out == "" {
		container.Io.Out = filepath.Join(dir, "-stdout")
	}
	if container.Io.Err == "" {
		container.Io.Err = filepath.Join(dir, "-stderr")
	}
	return container
}
