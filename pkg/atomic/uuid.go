package atomic

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gowebpki/jcs"
)

// NewUUID4 returns a new type 4 UUID.
func NewUUID4() string {
	return uuid.New().String()
}

// NewUUID5 returns a new type 5 UUID and panics if the provided UUID namespace is invalid.
func NewUUID5(namespace string, blob []byte) string {
	return uuid.NewSHA1(uuid.MustParse(namespace), blob).String()
}

// NewUUID5FromMap returns a new type 5 UUID from a map and panics if the provided UUID namespace is invalid.
func NewUUID5FromMap(namespace string, m map[string]interface{}) (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	cb, err := jcs.Transform(b)
	if err != nil {
		return "", err
	}
	return NewUUID5(namespace, cb), nil
}
