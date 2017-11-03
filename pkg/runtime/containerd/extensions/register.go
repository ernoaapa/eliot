package extensions

import (
	"strconv"

	"github.com/containerd/typeurl"
)

var VersionMajor = 1

func init() {
	const prefix = "types.can.io"
	// register TypeUrls for commonly marshaled external types
	major := strconv.Itoa(VersionMajor)
	typeurl.Register(&PipeSet{}, prefix, "containerd/extensions", major, "PipeSet")
}
