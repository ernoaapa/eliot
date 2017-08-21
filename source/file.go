package source

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ernoaapa/layeryd/model"
	"gopkg.in/yaml.v2"
)

type FileSource struct {
	filePath string
}

func NewFileSource(filePath string) *FileSource {
	return &FileSource{
		filePath,
	}
}

func (s *FileSource) GetState(info model.NodeInfo) (model.DesiredState, error) {
	data, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		return model.DesiredState{}, err
	}

	switch extension := filepath.Ext(s.filePath); extension {
	case ".yaml", ".yml":
		return unmarshalYaml(data)
	default:
		return model.DesiredState{}, fmt.Errorf("Invalid source file format: %s", extension)
	}
}

func unmarshalYaml(data []byte) (model.DesiredState, error) {
	state := &model.DesiredState{}
	err := yaml.Unmarshal(data, state)
	return *state, err
}
