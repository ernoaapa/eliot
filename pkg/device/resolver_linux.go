package device

import (
	"github.com/ernoaapa/layery/pkg/model"
	"github.com/matishsiao/goInfo"
)

// GetInfo resolves information about the device
func GetInfo(labels map[string]string) *model.DeviceInfo {
	osInfo := goInfo.GetInfo()

	return &model.DeviceInfo{
		Labels:   labels,
		Platform: osInfo.Platform,
		OS:       osInfo.GoOS,
		Kernel:   osInfo.Kernel,
		Core:     osInfo.Core,
		Hostname: osInfo.Hostname,
		CPUs:     osInfo.CPUs,

		MachineID: getInfoFromFiles([]string{
			"/etc/machine-id",
			"/var/lib/dbus/machine-id",
		}, failIfCannotResolve("MachineID")),

		SystemUUID: getInfoFromFiles([]string{
			"/sys/class/dmi/id/product_uuid",
			"/proc/device-tree/system-id",
			"/proc/device-tree/vm,uuid",
			"/etc/machine-id",
		}, failIfCannotResolve("SystemUUID")),

		BootID: getInfoFromFiles([]string{
			"/proc/sys/kernel/random/boot_id",
		}, failIfCannotResolve("BootID")),
	}
}
