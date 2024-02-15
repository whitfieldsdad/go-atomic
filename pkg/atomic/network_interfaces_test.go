package atomic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListNetworkInterfaces(t *testing.T) {
	nics, err := ListNetworkInterfaces()
	assert.Nil(t, err, "ListNetworkInterfaces() failed")
	assert.Greater(t, len(nics), 0, "Expected at least one network interface")
}

func TestGetNetworkInterface(t *testing.T) {
	nics, err := ListNetworkInterfaces()
	assert.Nil(t, err, "ListNetworkInterfaces() failed")

	for _, a := range nics {
		b, err := GetNetworkInterface(a.Name)
		assert.Nil(t, err, "GetNetworkInterface() failed")
		assert.Equal(t, a, *b, "Expected network interface names to match")
	}
}
