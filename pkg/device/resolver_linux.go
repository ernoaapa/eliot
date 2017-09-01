package device

import (
	"github.com/ernoaapa/can/pkg/model"
	"github.com/matishsiao/goInfo"
)

// GetInfo resolves information about the device
func (r *Resolver) GetInfo() *model.DeviceInfo {
	osInfo := goInfo.GetInfo()

	return &model.DeviceInfo{
		Labels:   r.labels,
		Platform: osInfo.Platform,
		OS:       osInfo.GoOS,
		Kernel:   osInfo.Kernel,
		Core:     osInfo.Core,
		Hostname: osInfo.Hostname,
		CPUs:     osInfo.CPUs,

		MachineID: resolveFirst(
			"MachineID",
			fromEnv("MACHINE_ID"),
			fromFiles([]string{
				"/etc/machine-id",
				"/var/lib/dbus/machine-id",
			}),
		),

		SystemUUID: resolveFirst(
			"SystemUUID",
			fromFiles([]string{
				"/sys/class/dmi/id/product_uuid",
				"/proc/device-tree/system-id",
				"/proc/device-tree/vm,uuid",
				"/etc/machine-id",
			}),
		),

		BootID: resolveFirst(
			"BootID",
			fromFiles([]string{
				"/proc/sys/kernel/random/boot_id",
			}),
		),
	}
}
