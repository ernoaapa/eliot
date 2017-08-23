package source

import "github.com/ernoaapa/layeryd/model"

type Source interface {
	// GetState(model.NodeInfo) (model.Pod, error)
	GetUpdates(model.NodeInfo) chan model.Pod
}
