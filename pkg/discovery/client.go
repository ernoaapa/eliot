package discovery

import (
	"context"
	"time"

	"github.com/ernoaapa/can/pkg/model"
	"github.com/grandcat/zeroconf"
	"github.com/pkg/errors"
)

// Devices search for devices in network for given timeout
func Devices(results chan<- model.DeviceInfo, timeout time.Duration) error {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return errors.Wrapf(err, "Failed to initialize new zeroconf resolver")
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(entries <-chan *zeroconf.ServiceEntry) {
		for entry := range entries {
			results <- model.DeviceInfo{
				Hostname:  entry.HostName,
				Addresses: append(entry.AddrIPv4, entry.AddrIPv6...),
				GrpcPort:  entry.Port,
			}
		}
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = resolver.Browse(ctx, ZeroConfServiceName, "", entries)
	if err != nil {
		return errors.Wrapf(err, "Failed to browse zeroconf devices")
	}

	<-ctx.Done()
	return nil
}
