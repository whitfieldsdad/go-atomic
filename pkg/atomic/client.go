package atomic

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gobwas/glob"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/afero/tarfs"
)

type Client struct {
	TaskDir         string `json:"task_dir"`
	TaskTemplateDir string `json:"task_template_dir"`
}

func NewClient() (*Client, error) {
	client := &Client{}
	return client, nil
}

func (c Client) ReadTasks(paths []string, q *TaskQuery) ([]Task, error) {
	var tasks []Task
	for _, path := range paths {
		taskptr, err := c.ReadTask(path)
		if err != nil {
			log.Errorf("Failed to read task: %s", err)
			continue
		}
		task := *taskptr
		if q != nil && !q.MatchesTask(task) {
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (c Client) ReadTask(path string) (*Task, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read task")
	}
	task, err := c.parseTask(b)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse task")
	}
	return task, nil
}

func (c Client) parseTask(b []byte) (*Task, error) {
	var task Task
	err := json.Unmarshal(b, &task)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal task")
	}
	if task.Id == "" {
		task.Id = NewUUID4()
	}
	stepTypeToStruct := map[StepType]interface{}{
		StepTypeExecuteCommand: ExecuteCommandStep{},
		StepTypeListProcesses:  ListProcessesStep{},
		StepTypeListUsers:      ListUsersStep{},
	}
	for i, step := range task.Steps {
		if step.Id == "" {
			task.Steps[i].Id = NewUUID4()
		}
		o, ok := stepTypeToStruct[step.Type]
		if !ok {
			return nil, errors.New("unknown step type")
		}
		err = mapstructure.Decode(step.Data, &o)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode step")
		}
		task.Steps[i].Data = o
	}
	return &task, nil
}

func (c Client) GenerateTaskTemplates(inputPath, outputDir string) error {
	templates, err := c.generateTaskTemplates(inputPath)
	if err != nil {
		return err
	}
	for _, template := range templates {
		filename := fmt.Sprintf("%s.json", template.Id)
		path := filepath.Join(outputDir, filename)
		err = writeJSONFile(path, template)
		if err != nil {
			return errors.Wrap(err, "failed to write task template")
		}
	}
	return nil
}

func (c Client) generateTaskTemplates(path string) ([]TaskTemplate, error) {
	path = getRealPath(path)
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if st.IsDir() {
		return c.generateTaskTemplatesFromDir(path)
	} else if strings.HasSuffix(path, ".tar.gz") {
		return c.generateTaskTemplatesFromTarball(path)
	} else if strings.HasSuffix(path, ".yaml") {
		bundle, err := readAtomicRedTeamYAMLFile(path)
		if err != nil {
			return nil, err
		}
		return bundle.GetTaskTemplates()
	}
	return nil, errors.New("unsupported source")
}

func (c Client) generateTaskTemplatesFromDir(path string) ([]TaskTemplate, error) {
	paths, err := findFiles(path, "T*/T*.yaml")
	if err != nil {
		return nil, err
	}
	var result []TaskTemplate
	for _, p := range paths {
		bundle, err := readAtomicRedTeamYAMLFile(p)
		if err != nil {
			log.Warnf("Failed to parse bundle: %s - %s", p, err)
			continue
		}
		templates, err := bundle.GetTaskTemplates()
		if err != nil {
			log.Warnf("Failed to get task templates: %s - %s", p, err)
			continue
		}
		result = append(result, templates...)
	}
	return result, nil
}

func (c Client) generateTaskTemplatesFromTarball(path string) ([]TaskTemplate, error) {
	tf, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer tf.Close()

	gr, err := gzip.NewReader(tf)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	tfs := tarfs.New(tar.NewReader(gr))
	afs := &afero.Afero{Fs: tfs}

	var paths []string
	afs.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		g := glob.MustCompile("*/T*/T*.yaml")
		if g.Match(path) {
			paths = append(paths, path)
		}
		return nil
	})
	var allTemplates []TaskTemplate
	for _, path := range paths {
		f, err := afs.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		b, err := afero.ReadAll(f)
		if err != nil {
			return nil, err
		}
		bundle, err := parseAtomicRedTeamYAML(b)
		if err != nil {
			log.Warnf("Failed to parse bundle: %s - %s", path, err)
		}
		templates, err := bundle.GetTaskTemplates()
		if err != nil {
			log.Warnf("Failed to get task templates: %s - %s", path, err)
		}
		allTemplates = append(allTemplates, templates...)
	}
	return allTemplates, nil
}
