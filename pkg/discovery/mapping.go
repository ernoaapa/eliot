package discovery

import (
	"strings"

	"github.com/ernoaapa/eliot/pkg/model"
	"github.com/grandcat/zeroconf"
)

// MapToInternalModel takes zeroconf entry and maps it to internal DeviceInfo model
func MapToInternalModel(entry *zeroconf.ServiceEntry) model.DeviceInfo {
	version := "unknown"

	for _, val := range entry.Text {
		parts := strings.SplitN(val, "=", 2)
		if len(parts) == 2 && parts[0] == "v" {
			version = parts[1]
		}
	}

	return model.DeviceInfo{
		Hostname:  entry.HostName,
		Addresses: append(entry.AddrIPv4, entry.AddrIPv6...),
		GrpcPort:  entry.Port,
		Version:   version,
	}
}
