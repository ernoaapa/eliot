package build

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/ernoaapa/eliot/pkg/cmd"
	"github.com/pkg/errors"
)

// ResolveLinuxkitConfig resolves and reads linuxkit config from given source.
// Source can be url or path to the file.
func ResolveLinuxkitConfig(source string) (linuxkit []byte, err error) {

	if source == "" {
		// Default to default rpi3 Linuxkit config
		source = "https://raw.githubusercontent.com/ernoaapa/eliot-os/master/rpi3.yml"
	}

	if isValidFile(source) {
		linuxkit, err = ioutil.ReadFile(source)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to read Linuxkit file")
		}
	} else if isValidURL(source) {
		linuxkit, err = getContent(source)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to fetch Linuxkit config")
		}
	} else if cmd.IsPipingIn() {
		linuxkit, err = cmd.ReadAllStdin()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to read Linuxkit config from stdin")
		}
	} else {
		return nil, errors.New("No Linuxkit config defined")
	}

	if len(linuxkit) == 0 {
		return nil, errors.New("Invalid Linuxkit config")
	}

	return linuxkit, nil
}

// BuildImage builds given Linuxkit config and returns the image as io.ReadClose or error
// Note: you must call image.Close()
func BuildImage(serverURL, outputType, outputFormat string, config []byte) (io.ReadCloser, error) {
	res, err := http.Post(fmt.Sprintf("%s/linuxkit/%s/build/%s?output=%s", serverURL, "eli-cli", outputType, outputFormat), "application/yml", bytes.NewReader(config))
	if err != nil {
		return nil, errors.Wrap(err, "Error while making request to Linuxkit build server")
	}
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return res.Body, nil
	}

	errmsg, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read build error response")
	}
	return nil, fmt.Errorf("Build failed: %s", errmsg)
}

// isValidURL tests a string to determine if it is a url or not.
func isValidURL(toTest string) bool {
	if _, err := url.ParseRequestURI(toTest); err != nil {
		return false
	}
	return true
}

// getContent fetch url and returns all content
func getContent(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// isValidFile
func isValidFile(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}
