package discovery

import (
	"net"
	"testing"

	"github.com/grandcat/zeroconf"
	"github.com/stretchr/testify/assert"
)

func TestMapToInternalModel(t *testing.T) {
	result := MapToInternalModel(&zeroconf.ServiceEntry{
		HostName: "hostname",
		AddrIPv4: []net.IP{net.IPv4zero},
		AddrIPv6: []net.IP{net.IPv6loopback},
		Text:     []string{"v=1.2.3-abcd"},
	})

	assert.Equal(t, "hostname", result.Hostname)
	assert.Equal(t, "1.2.3-abcd", result.Version)
	assert.Equal(t, []net.IP{net.IPv4zero, net.IPv6loopback}, result.Addresses)
}
