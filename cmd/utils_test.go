package cmd

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestFileExists(t *testing.T) {
	file, err := ioutil.TempFile("", "TestFileExists")
	assert.NoError(t, err, "Failed to create temporary file for test")
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	assert.True(t, fileExists(file.Name()), "should return true if file exists")
	assert.False(t, fileExists("/this/file/dont/exist"), "should return false if file does not exists")
}

func TestGetNoLabels(t *testing.T) {
	flags := flag.NewFlagSet("test", 0)
	flags.String("labels", "", "")

	clicontext := cli.NewContext(nil, flags, nil)

	labels := GetLabels(clicontext)

	assert.Equal(t, map[string]string{}, labels)
}

func TestGetSingleLabel(t *testing.T) {
	flags := flag.NewFlagSet("test", 0)
	flags.String("labels", "", "")

	clicontext := cli.NewContext(nil, flags, nil)
	flags.Parse([]string{"--labels", "foo=bar"})

	labels := GetLabels(clicontext)

	assert.Equal(t, map[string]string{
		"foo": "bar",
	}, labels)
}

func TestGetMultipleLabels(t *testing.T) {
	flags := flag.NewFlagSet("test", 0)
	flags.String("labels", "", "")

	flags.Parse([]string{"--labels", "foo=bar,doo=daa,ugh=12.3.4"})
	clicontext := cli.NewContext(nil, flags, nil)

	labels := GetLabels(clicontext)

	assert.Equal(t, map[string]string{
		"foo": "bar",
		"doo": "daa",
		"ugh": "12.3.4",
	}, labels)
}

func TestParseMountFlag(t *testing.T) {
	result, err := parseMountFlag("type=foo,source=/path,destination=/target,options=rbind:rw")
	assert.NoError(t, err)
	assert.Equal(t, "foo", result.Type)
	assert.Equal(t, "/path", result.Source)
	assert.Equal(t, "/target", result.Destination)
	assert.Equal(t, []string{"rbind", "rw"}, result.Options)
}

func TestParseBindFlag(t *testing.T) {
	result, err := parseBindFlag("/source:/target:rshared")
	assert.NoError(t, err)
	assert.Equal(t, "bind", result.Type)
	assert.Equal(t, "/source", result.Source)
	assert.Equal(t, "/target", result.Destination)
	assert.Equal(t, []string{"rshared", "rbind"}, result.Options)
}

func TestExpandToFQIN(t *testing.T) {
	assert.Equal(t, "docker.io/eaapa/hello-world:latest", ExpandToFQIN("eaapa/hello-world"))
	assert.Equal(t, "otherhost.io/eaapa/hello-world:latest", ExpandToFQIN("otherhost.io/eaapa/hello-world"))
	assert.Equal(t, "docker.io/eaapa/hello-world:latest", ExpandToFQIN("eaapa/hello-world:latest"))
	assert.Equal(t, "docker.io/library/nginx:tag1", ExpandToFQIN("nginx:tag1"))
	assert.Equal(t, "docker.io/library/nginx:latest", ExpandToFQIN("nginx"))
}
