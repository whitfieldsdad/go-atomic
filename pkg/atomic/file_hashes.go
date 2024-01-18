package atomic

import "time"

type FileHashes struct {
	Id     string    `json:"id"`
	Time   time.Time `json:"time"`
	Path   string    `json:"path"`
	Hashes Hashes    `json:"hashes"`
}
