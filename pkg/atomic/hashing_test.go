package atomic

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHashes(t *testing.T) {
	b := []byte("Hello, world!")
	r := bytes.NewReader(b)
	hashes, err := GetHashes(r)
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, hashes, "Hashes should not be nil")

	expected := &Hashes{
		MD5:    "6cd3556deb0da54bca060b4c39479839",
		SHA1:   "943a702d06f34599aee1f8da8ef9f7296031d699",
		SHA256: "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3",
	}
	assert.Equal(t, expected, hashes, "Hashes should match")
}

