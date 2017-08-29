package utils

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
