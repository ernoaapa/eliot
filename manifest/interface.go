package manifest

import "github.com/ernoaapa/layeryd/model"

// Source is interface for all state sources
type Source interface {
	GetUpdates() chan []model.Pod
}
