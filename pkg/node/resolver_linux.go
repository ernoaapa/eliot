package node

import (
	"os"
	"runtime"
	"strings"
	"syscall"

	"github.com/ernoaapa/eliot/pkg/model"
	log "github.com/sirupsen/logrus"
)

var mountTableFile = "/etc/mtab"

// GetInfo resolves information about the node
func (r *Resolver) GetInfo() *model.NodeInfo {
	hostname, _ := os.Hostname()
	return &model.NodeInfo{
		Version:   r.version,
		Labels:    r.labels,
		Arch:      runtime.GOARCH,
		OS:        runtime.GOOS,
		Hostname:  hostname,
		Addresses: getAddresses(),
		GrpcPort:  r.grpcPort,

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
		Filesystems: resolveFilesystems(),
	}
}

func resolveUptime() uint64 {
	sysinfo := syscall.Sysinfo_t{}

	if err := syscall.Sysinfo(&sysinfo); err != nil {
		return err
	}

	return sysinfo.Uptime
}

// resolveFilesystems resolves filesystems from /etc/mtab file
func resolveFilesystems() []model.Filesystem {
	result := []model.Filesystem{}

	err := readFile(mountTableFile, func(line string) error {
		fields := strings.Fields(line)

		devName := fields[0]
		dirName := fields[1]
		sysTypeName := fields[2]

		total, free, available, err := getFilesystemUsage(dirName)
		if err != nil {
			return err
		}

		result = append(result, model.Filesystem{
			Filesystem: devName,
			TypeName:   sysTypeName,
			MountDir:   dirName,
			Total:      total,
			Free:       free,
			Available:  available,
		})

		return nil
	})

	if err != nil {
		log.Errorf("Failed to resolve filesystems from %s, fallback to empty list. Error: %s", mountTableFile, err)
	}
	return result
}

func getFilesystemUsage(path string) (total uint64, free uint64, available uint64, err error) {
	stat := syscall.Statfs_t{}
	err = syscall.Statfs(path, &stat)
	if err != nil {
		return 0, 0, 0, err
	}

	return uint64(stat.Blocks) * uint64(stat.Bsize), uint64(stat.Bfree) * uint64(stat.Bsize), uint64(stat.Bavail) * uint64(stat.Bsize), nil
}
