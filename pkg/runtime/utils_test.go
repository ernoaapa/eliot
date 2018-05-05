package runtime

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ernoaapa/eliot/pkg/fs"
	"github.com/ernoaapa/eliot/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestEnsureMountSourceDirExistsCreatesDirectory(t *testing.T) {
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	source := filepath.Join(dir, "temp", "directory")
	assert.False(t, fs.DirExist(source))
	assert.NoError(t, ensureMountSourceDirExists([]model.Mount{{Source: source}}))
	assert.True(t, fs.DirExist(source))
}

func TestEnsureMountSourceDirExistsSkipFiles(t *testing.T) {
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	source := filepath.Join(dir, "file")
	assert.NoError(t, ioutil.WriteFile(source, []byte("foobar"), os.ModePerm))

	assert.NoError(t, ensureMountSourceDirExists([]model.Mount{{Source: source}}))
}

func TestGetValues(t *testing.T) {
	expected := &model.Pod{}
	result := getValues(map[string]*model.Pod{
		"foo": expected,
	})

	assert.Equal(t, []model.Pod{*expected}, result)
}
