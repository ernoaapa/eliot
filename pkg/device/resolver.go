package device

import (
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Resolver provides information about the device
type Resolver struct {
	labels map[string]string
}

// NewResolver creates new resolver with static device labels
func NewResolver(labels map[string]string) *Resolver {
	return &Resolver{
		labels: withHostLabels(labels),
	}
}

func withHostLabels(labels map[string]string) map[string]string {
	if _, exist := labels["arch"]; !exist {
		labels["arch"] = runtime.GOARCH
	}

	if _, exist := labels["os"]; !exist {
		labels["os"] = runtime.GOOS
	}
	return labels
}

func getAddresses() (addresses []net.IP) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Errorf("Unable to resolve network interfaces, cannot expose addresses information: %s", err)
		return addresses
	}

	if len(ifaces) == 0 {
		log.Warn("Were not able to resolve any network interface")
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Errorf("Error while resolving interface [%s] addresses: %s", iface.Name, err)
			continue
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
			case *net.IPAddr:
				addresses = append(addresses, v.IP)
			}
		}
	}
	return addresses
}

func resolveFirst(name string, resolvers ...func() string) string {
	for _, resolver := range resolvers {
		result := resolver()
		if result != "" {
			return result
		}
	}

	log.Fatalf("Failed to resolve %s no default provided!", name)
	return ""
}

func fromEnv(name string) func() string {
	return func() string {
		return os.Getenv(name)
	}
}

func fromFiles(filePaths []string) func() string {
	return func() string {
		for _, file := range filePaths {
			info, err := ioutil.ReadFile(file)
			if err == nil {
				return strings.TrimSpace(string(info))
			}
		}
		return ""
	}
}

func static(value string) func() string {
	return func() string {
		return value
	}
}
