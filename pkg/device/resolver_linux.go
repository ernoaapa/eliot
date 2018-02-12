package device

import (
	"os"
	"runtime"

	"github.com/ernoaapa/eliot/pkg/model"
)

// GetInfo resolves information about the device
func (r *Resolver) GetInfo(grpcPort int, version string) *model.DeviceInfo {
	hostname, _ := os.Hostname()
	return &model.DeviceInfo{
		Version:   version,
		Labels:    r.labels,
		Arch:      runtime.GOARCH,
		OS:        runtime.GOOS,
		Hostname:  hostname,
		Addresses: getAddresses(),
		GrpcPort:  grpcPort,

		MachineID: resolveFirst(
			"MachineID",
			fromEnv("MACHINE_ID"),
			fromFiles([]string{
				"/etc/machine-id",
				"/var/lib/dbus/machine-id",
			}),
			static("unknown"),
		),

		SystemUUID: resolveFirst(
			"SystemUUID",
			fromFiles([]string{
				"/sys/class/dmi/id/product_uuid",
				"/proc/device-tree/system-id",
				"/proc/device-tree/vm,uuid",
				"/etc/machine-id",
			}),
			static("unknown"),
		),

		BootID: resolveFirst(
			"BootID",
			fromFiles([]string{
				"/proc/sys/kernel/random/boot_id",
			}),
			static("unknown"),
		),
	}
}
