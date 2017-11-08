package extensions

import (
	"strconv"

	"github.com/containerd/typeurl"
)

var versionMajor = 1

func init() {
	const prefix = "types.can.io"
	// register TypeUrls for commonly marshaled external types
	major := strconv.Itoa(versionMajor)
	typeurl.Register(&PipeSet{}, prefix, "containerd/extensions", major, "PipeSet")
	typeurl.Register(&ContainerLifecycle{}, prefix, "containerd/extensions", major, "ContainerLifecycle")
}
