package status

import "github.com/ernoaapa/layeryd/model"

// Reporter sends information about current status
type Reporter interface {
	Report(model.DeviceInfo, model.DeviceState) error
}
