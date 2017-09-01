package manifest

import "github.com/ernoaapa/can/pkg/model"

// Source is interface for all state sources
type Source interface {
	GetUpdates() chan []model.Pod
}
