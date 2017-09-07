package manifest

import (
	"encoding/json"

	"github.com/ernoaapa/can/pkg/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

func unmarshalYaml(data []byte) ([]model.Pod, error) {
	target := &[]model.Pod{}

	unmarshalErr := yaml.Unmarshal(data, target)
	if unmarshalErr != nil {
		return []model.Pod{}, errors.Wrapf(unmarshalErr, "Unable to parse Yaml data")
	}

	pods := model.Defaults(*target)

	validationErr := model.Validate(pods)
	if validationErr != nil {
		return pods, errors.Wrapf(validationErr, "Invalid pod definitions")
	}

	return pods, nil
}

func unmarshalJSON(data []byte) ([]model.Pod, error) {
	target := &[]model.Pod{}

	unmarshalErr := json.Unmarshal(data, target)
	if unmarshalErr != nil {
		log.Debugf("Unable to parse JSON: %s", string(data[:]))
		return []model.Pod{}, errors.Wrapf(unmarshalErr, "Unable to parse JSON data")
	}

	pods := model.Defaults(*target)

	validationErr := model.Validate(pods)
	if validationErr != nil {
		return pods, errors.Wrapf(validationErr, "Invalid pod definitions")
	}

	return pods, nil
}
