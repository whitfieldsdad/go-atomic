package atomic

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

type FileHashes struct {
	Artifact `json:",inline"`
	Path     string `json:"path"`
	Hashes   Hashes `json:"hashes"`
}

type File struct {
	Artifact
	Path     string  `json:"path"`
	Filename string  `json:"filename"`
	Hashes   *Hashes `json:"hashes"`
}

func NewFile(path string) File {
	filename := filepath.Base(path)
	return File{
		Artifact: NewArtifact(),
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

func getRealPath(path string) string {
	if strings.Contains(path, "~") {
		homedir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to lookup home directory: %s", err)
		}
		path = strings.ReplaceAll(path, "~", homedir)
	}
	return path
}

func findFiles(root, pathPattern string) ([]string, error) {
	var paths []string

	root = getRealPath(root)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		relativePath, _ := filepath.Rel(root, path)
		if pathPattern != "" {
			ok, _ := filepath.Match(pathPattern, relativePath)
			if !ok {
				return nil
			}
		}
		paths = append(paths, path)
		return nil
	})
	return paths, err
}

func writeJSONFile(path string, v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}
