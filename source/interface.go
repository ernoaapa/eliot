package source

import "github.com/ernoaapa/layeryd/model"

type Source interface {
	GetState(model.NodeInfo) (model.DesiredState, error)
}
