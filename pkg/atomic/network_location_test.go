package atomic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNetworkLocation(t *testing.T) {
	loc, err := GetNetworkLocation()
	assert.Nil(t, err)
	assert.NotNil(t, loc)
	assert.NotEmpty(t, loc.Hostname)
	assert.NotEmpty(t, loc.MACAddress)

	ips := append(loc.IPv4Addresses, loc.IPv6Addresses...)
	assert.NotEmpty(t, ips)
}
