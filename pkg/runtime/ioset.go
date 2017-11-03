package runtime

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// IOSet represents container process stdin,stdout,stderr files
type IOSet struct {
	Stdin  string
	Stdout string
	Stderr string
}

// NewIOSet creates new unique IOSet for container
func NewIOSet(id string) (*IOSet, error) {
	root := "/run/containerd/fifo"
	if err := os.MkdirAll(root, 0700); err != nil {
		return nil, err
	}
	dir, err := ioutil.TempDir(root, "")
	if err != nil {
		return nil, err
	}
	return &IOSet{
		Stdin:  filepath.Join(dir, id+"-stdin"),
		Stdout: filepath.Join(dir, id+"-stdout"),
		Stderr: filepath.Join(dir, id+"-stderr"),
	}, nil
}

// PipeStdoutTo updates the IOSet stdout to another IOSet stdin
func (s *IOSet) PipeStdoutTo(target *IOSet) {
	s.Stdout = target.Stdin
}
