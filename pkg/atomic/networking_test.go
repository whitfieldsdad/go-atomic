package atomic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostname(t *testing.T) {
	hostname, err := GetHostname()
	assert.Nil(t, err)
	assert.NotEmpty(t, hostname)
}

func TestGetPrimaryIPv4Address(t *testing.T) {
	ip, err := GetPrimaryIPv4Address()
	assert.Nil(t, err)
	assert.NotEmpty(t, ip)
}

func TestIsLoopbackAddress(t *testing.T) {
	assert.True(t, IsLoopbackAddress("127.0.0.1"))
}

func TestGetHostnameFromIPv4Address(t *testing.T) {
	ip, err := GetPrimaryIPv4Address()
	assert.Nil(t, err)
	assert.NotEmpty(t, ip)

	hostname, err := GetHostnameFromIPAddress(ip)
	assert.Nil(t, err)
	assert.NotEmpty(t, hostname)
}
