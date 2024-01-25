package atomic

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

type Client struct {
	AtomicsPath string `json:"atomics_path"`
}

func NewClient(atomicsPath string) (*Client, error) {
	if atomicsPath == "" {
		return nil, errors.New("Path to atomics folder is required")
	}
	client := &Client{
		AtomicsPath: atomicsPath,
	}
	return client, nil
}

func (c Client) ListTasks(q *TaskQuery) ([]Task, error) {
	tasks, err := c.listAtomicRedTeamTasks()
	if err != nil {
		return nil, err
	}
	if q != nil {
		var filteredTasks []Task
		for _, task := range tasks {
			if q.Matches(task) {
				filteredTasks = append(filteredTasks, task)
			}
		}
		tasks = filteredTasks
	}
	return tasks, nil
}

func (c Client) listAtomicRedTeamTasks() ([]Task, error) {
	paths, err := findFilePaths(c.AtomicsPath, "T*/T*.yaml")
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for _, path := range paths {
		bundle, err := ReadAtomicRedTeamYAMLFile(path)
		if err != nil {
			log.Warnf("Failed to parse bundle: %s - %s", path, err)
			continue
		}
		ts, err := bundle.GetTasks()
		if err != nil {
			continue
		}
		tasks = append(tasks, ts...)
	}
	return tasks, nil
}

func getRealPath(path string) (string, error) {
	if strings.Contains(path, "~") {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = strings.ReplaceAll(path, "~", homedir)
	}
	return path, nil
}

func findFilePaths(root, pathPattern string) ([]string, error) {
	var paths []string

	root, err := getRealPath(root)
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			log.Fatalf("%s is nil", path)
		}
		if info.IsDir() {
			return nil
		}
		relativePath, _ := filepath.Rel(root, path)
		if pathPattern != "" {
			ok, err := filepath.Match(pathPattern, relativePath)
			if err != nil {
				log.Fatalf("Failed to compare %s and %s", path, pathPattern)
			}
			if !ok {
				return nil
			}
		}
		paths = append(paths, path)
		return nil
	})
	return paths, err
}
