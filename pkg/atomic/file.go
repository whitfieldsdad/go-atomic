package atomic

import (
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type File struct {
	Id       string    `json:"id"`
	Time     time.Time `json:"time"`
	Path     string    `json:"path"`
	Filename string    `json:"filename"`
	Hashes   *Hashes   `json:"hashes"`
}

func NewFile(path string) File {
	filename := filepath.Base(path)
	return File{
		Id:       uuid.NewString(),
		Time:     time.Now(),
		Filename: filename,
		Path:     path,
	}
}

func GetFile(path string) (*File, error) {
	file := NewFile(path)
	hashes, err := GetFileHashes(file.Path)
	if err != nil {
		return nil, err
	}
	file.Hashes = hashes
	return &file, nil
}
