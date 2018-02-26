package node

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveFirst(t *testing.T) {
	assert.Equal(t, "first", resolveFirst("foofield", func() string {
		return "first"
	}, func() string {
		return "second"
	}))

	assert.Equal(t, "second", resolveFirst("foofield", func() string {
		return "" // Mimic the case that cannot resolve
	}, func() string {
		return "second"
	}))
}

func TestFromEnv(t *testing.T) {
	os.Setenv("TESTING", "foo")
	assert.Equal(t, "foo", fromEnv("TESTING")())
	assert.Equal(t, "", fromEnv("DONT_EXIST")())
}

func TestFromFiles(t *testing.T) {
	dir, createErr := ioutil.TempDir("", "example")
	assert.NoError(t, createErr, "Unable to create temp file")
	filePath := fmt.Sprintf("%s/%s", dir, "test.yml")

	writeErr := ioutil.WriteFile(filePath, []byte("foobar"), 0666)
	assert.NoError(t, writeErr, "Unable to write to temporary file")

	assert.Equal(t, "foobar", fromFiles([]string{filePath})())

}
