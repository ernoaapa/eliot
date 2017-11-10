package discovery

import (
	"context"
	"time"

	"github.com/ernoaapa/elliot/pkg/model"
	"github.com/grandcat/zeroconf"
	"github.com/pkg/errors"
)

// Devices return list of DeviceInfos synchronously with given timeout
func Devices(timeout time.Duration) (devices []model.DeviceInfo, err error) {
	results := make(chan model.DeviceInfo)
	defer close(results)

	go func() {
		for device := range results {
			devices = append(devices, device)
		}
	}()

	err = DevicesAsync(results, timeout)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// DevicesAsync search for devices in network asynchronously for given timeout
func DevicesAsync(results chan<- model.DeviceInfo, timeout time.Duration) error {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return errors.Wrapf(err, "Failed to initialize new zeroconf resolver")
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(entries <-chan *zeroconf.ServiceEntry) {
		for entry := range entries {
			results <- MapToInternalModel(entry)
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
