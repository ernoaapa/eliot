package resolve

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"

	pods "github.com/ernoaapa/elliot/pkg/api/services/pods/v1"
	"github.com/ernoaapa/elliot/pkg/fs"
	"github.com/pkg/errors"
)

// Pods resolve list of Pod resources
// source can be
// - directory of yaml specs
// - yaml spec file
// - url to download yaml spec
func Pods(sources []string) (result []*pods.Pod, err error) {
	for _, source := range sources {
		if fs.FileExist(source) {
			resources, err := readFileSource(source)
			if err != nil {
				return result, errors.Wrapf(err, "Failed to read pod spec file %s", source)
			}
			result = append(result, resources...)
		} else if fs.DirExist(source) {
			files, err := ioutil.ReadDir(source)
			if err != nil {
				return result, errors.Wrapf(err, "Failed to read pod spec directory %s", source)
			}
			for _, file := range files {
				if !file.IsDir() {
					resources, err := readFileSource(filepath.Join(source, file.Name()))
					if err != nil {
						return result, errors.Wrapf(err, "Failed to read pod spec file %s", source)
					}
					result = append(result, resources...)
				}
			}
		} else if validURL(source) {
			response, err := http.Get(source)
			if err != nil {
				return result, errors.Wrapf(err, "Failed to load spec from url: %s", source)
			}
			defer response.Body.Close()

			data, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return result, errors.Wrapf(err, "Failed to get spec response from url: %s", source)
			}
			resources, err := pods.UnmarshalYaml(data)
			if err != nil {
				return result, errors.Wrapf(err, "Failed to read pod spec response from url: %s", source)
			}

			result = append(result, resources...)
		} else {
			return result, fmt.Errorf("Unknown source %s. Must be file, directory or url", source)
		}
	}
	return result, nil
}

func readFileSource(path string) ([]*pods.Pod, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return []*pods.Pod{}, errors.Wrapf(err, "Failed to read pod spec file %s", path)
	}

	return pods.UnmarshalYaml(data)
}

func validURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}
