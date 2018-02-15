package cmd

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	containers "github.com/ernoaapa/eliot/pkg/api/services/containers/v1"
	"github.com/ernoaapa/eliot/pkg/config"
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
	result, err := ParseBindFlag("/source:/target:rshared")
	assert.NoError(t, err)
	assert.Equal(t, "bind", result.Type)
	assert.Equal(t, "/source", result.Source)
	assert.Equal(t, "/target", result.Destination)
	assert.Equal(t, []string{"rshared", "rbind"}, result.Options)
}

func TestGetCurrentDirectory(t *testing.T) {
	assert.NotEmpty(t, GetCurrentDirectory())
}

func TestDropDoubleDash(t *testing.T) {
	assert.Equal(t, []string{"bash", "-il"}, DropDoubleDash([]string{"--", "bash", "-il"}))
	assert.Equal(t, []string{"bash", "-il"}, DropDoubleDash([]string{"foo", "--", "bash", "-il"}))
	assert.Equal(t, []string{"bash", "-il"}, DropDoubleDash([]string{"bash", "-il"}))
	assert.Equal(t, []string{}, DropDoubleDash([]string{"--"}))
	assert.Equal(t, []string{}, DropDoubleDash([]string{}))
}

func TestResolveSingleContainerID(t *testing.T) {
	containerID, err := ResolveContainerID([]*containers.ContainerStatus{
		{ContainerID: "1", Name: "foo"},
	}, "")
	assert.NoError(t, err)
	assert.Equal(t, "1", containerID)
}

func TestMultiResolveContainerID(t *testing.T) {
	containerID, err := ResolveContainerID([]*containers.ContainerStatus{
		{ContainerID: "1", Name: "foo"},
		{ContainerID: "2", Name: "bar"},
	}, "bar")
	assert.NoError(t, err)
	assert.Equal(t, "2", containerID)
}

func TestMultiResolveContainerIDFail(t *testing.T) {
	_, err := ResolveContainerID([]*containers.ContainerStatus{
		{ContainerID: "1", Name: "foo"},
		{ContainerID: "2", Name: "bar"},
	}, "")
	assert.Error(t, err)
}

func TestEmptyResolveContainerIDFail(t *testing.T) {
	_, err := ResolveContainerID([]*containers.ContainerStatus{}, "")
	assert.Error(t, err)
}

func TestGetConfigProviderWithEndpointFlag(t *testing.T) {
	flags := flag.NewFlagSet("test", 0)
	flags.String("endpoint", "", "")

	flags.Parse([]string{"--endpoint", "1.2.3.4:5000"})
	clicontext := cli.NewContext(nil, flags, nil)

	provider := GetConfigProvider(clicontext)

	assert.Equal(t, []config.Endpoint{{
		Name: "1.2.3.4:5000",
		URL:  "1.2.3.4:5000",
	}}, provider.GetEndpoints(), "")
}
