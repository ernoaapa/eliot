package pods

import (
	"bufio"
	"bytes"
	"encoding/json"

	utils "github.com/ernoaapa/can/pkg/utils/yaml"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// UnmarshalYaml reads v1 Pods data in YAML format and unmarshals it to v1 api model
func UnmarshalYaml(data []byte) ([]*Pod, error) {
	result := []*Pod{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(utils.SplitYAMLDocument)

	for scanner.Scan() {
		target := &Pod{}
		unmarshalErr := yaml.Unmarshal(scanner.Bytes(), target)
		if unmarshalErr != nil {
			return result, errors.Wrapf(unmarshalErr, "Unable to parse Yaml data")
		}
		result = append(result, target)
	}

	return Defaults(result), nil
}

// UnmarshalListYaml reads list of v1 Pods data in YAML format and unmarshals it to v1 api model
func UnmarshalListYaml(data []byte) ([]*Pod, error) {
	target := &[]*Pod{}

	unmarshalErr := yaml.Unmarshal(data, target)
	if unmarshalErr != nil {
		return []*Pod{}, errors.Wrapf(unmarshalErr, "Unable to parse Yaml data")
	}

	return Defaults(*target), nil
}

// UnmarshalListJSON reads v1 Pods data in JSON format and unmarshals it to v1 api model
func UnmarshalListJSON(data []byte) ([]*Pod, error) {
	target := &[]*Pod{}

	unmarshalErr := json.Unmarshal(data, target)
	if unmarshalErr != nil {
		log.Debugf("Unable to parse JSON: %s", string(data[:]))
		return []*Pod{}, errors.Wrapf(unmarshalErr, "Unable to parse JSON data")
	}

	return Defaults(*target), nil
}
