package manifest

import (
	"github.com/ernoaapa/layery/pkg/model"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func unmarshalYaml(data []byte) ([]model.Pod, error) {
	target := &[]model.Pod{}

	unmarshalErr := yaml.Unmarshal(data, target)
	if unmarshalErr != nil {
		return []model.Pod{}, errors.Wrapf(unmarshalErr, "Unable to read yaml file")
	}

	pods := model.Defaults(*target)

	validationErr := model.Validate(pods)
	if validationErr != nil {
		return pods, errors.Wrapf(validationErr, "Invalid pod definitions")
	}

	return pods, nil
}
