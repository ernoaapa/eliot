package discovery

import (
	"net"
	"strings"

	node "github.com/ernoaapa/eliot/pkg/api/services/node/v1"
	"github.com/grandcat/zeroconf"
)

func MapToAPIModel(entry *zeroconf.ServiceEntry) *node.Info {
	version := "unknown"

	for _, val := range entry.Text {
		parts := strings.SplitN(val, "=", 2)
		if len(parts) == 2 && parts[0] == "v" {
			version = parts[1]
		}
	}

	return &node.Info{
		Hostname:  entry.HostName,
		Addresses: addressesToString(append(entry.AddrIPv4, entry.AddrIPv6...)),
		GrpcPort:  int64(entry.Port),
		Version:   version,
	}
}

func addressesToString(addresses []net.IP) (result []string) {
	for _, ip := range addresses {
		result = append(result, ip.String())
	}
	return result
}
