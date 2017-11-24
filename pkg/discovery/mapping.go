package discovery

import (
	"github.com/ernoaapa/eliot/pkg/model"
	"github.com/grandcat/zeroconf"
)

func MapToInternalModel(entry *zeroconf.ServiceEntry) model.DeviceInfo {
	return model.DeviceInfo{
		Hostname:  entry.HostName,
		Addresses: append(entry.AddrIPv4, entry.AddrIPv6...),
		GrpcPort:  entry.Port,
	}
}
