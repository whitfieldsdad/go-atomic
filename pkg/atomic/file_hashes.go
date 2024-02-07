package atomic

type FileHashes struct {
	Artifact `json:",inline"`
	Path     string `json:"path"`
	Hashes   Hashes `json:"hashes"`
}
