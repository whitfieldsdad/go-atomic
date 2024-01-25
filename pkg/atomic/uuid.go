package atomic

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gowebpki/jcs"
)

func NewUUID4() string {
	return uuid.New().String()
}

func NewUUID5(namespace string, blob []byte) string {
	return uuid.NewSHA1(uuid.MustParse(namespace), blob).String()
}

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
