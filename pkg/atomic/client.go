package atomic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

var (
	DefaultWorkdir = "./workdir"
)

type ClientConfig struct {
	WorkdirConfig WorkdirConfig `json:"workdir_config"`
}

type WorkdirConfig struct {
	TaskDir                 string `json:"task_dir"`
	TaskTemplateDir         string `json:"task_template_dir"`
	TaskInvocationDir       string `json:"task_invocation_dir"`
	TaskInvocationStatusDir string `json:"task_invocation_status_dir"`
	TaskInvocationResultDir string `json:"task_invocation_result_dir"`
	ArtifactDir             string `json:"artifact_dir"`
}

type Client struct {
	Config ClientConfig
}

func NewClient(workdir string) (*Client, error) {
	if workdir == "" {
		workdir = DefaultWorkdir
	}
	client := &Client{
		Config: ClientConfig{
			WorkdirConfig: getWorkdirConfig(workdir),
		},
	}
	return client, nil
}

func getWorkdirConfig(workdir string) WorkdirConfig {
	workdir = getRealPath(workdir)
	return WorkdirConfig{
		TaskDir:                 filepath.Join(workdir, "tasks"),
		TaskTemplateDir:         filepath.Join(workdir, "task_templates"),
		TaskInvocationDir:       filepath.Join(workdir, "task_invocations"),
		TaskInvocationStatusDir: filepath.Join(workdir, "task_invocation_status"),
		TaskInvocationResultDir: filepath.Join(workdir, "task_invocation_results"),
		ArtifactDir:             filepath.Join(workdir, "artifacts"),
	}
}

func (c Client) ReadTasks(paths []string) ([]Task, error) {
	var tasks []Task
	for _, p := range paths {
		task, err := c.ReadTask(p)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func (c Client) ReadTask(path string) (*Task, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var task Task
	err = json.NewDecoder(file).Decode(&task)
	if err != nil {
		return nil, err
	}
	if task.Id == "" {
		task.Id = NewUUID4()
	}
	stepTypeToStruct := map[StepType]interface{}{
		StepTypeExecuteCommand: ExecuteCommandStep{},
		StepTypeListProcesses:  ListProcessesStep{},
	}
	for i, step := range task.Steps {
		if step.Id == "" {
			task.Steps[i].Id = NewUUID4()
		}
		o, ok := stepTypeToStruct[step.Type]
		if !ok {
			return nil, fmt.Errorf("unknown step type: %s", step.Type)
		}
		err = mapstructure.Decode(step.Data, &o)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode step")
		}
		task.Steps[i].Data = o
	}
	return &task, nil
}

func (c Client) ListTasks(q *TaskQuery) ([]Task, error) {
	paths, err := findFiles(c.Config.WorkdirConfig.TaskDir, "*.json")
	if err != nil {
		return nil, err
	}
	var tasks []Task
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		var task Task
		err = json.NewDecoder(file).Decode(&task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (c Client) ListTaskTemplates(q *TaskQuery) ([]TaskTemplate, error) {
	paths, err := findFiles(c.Config.WorkdirConfig.TaskTemplateDir, "*.json")
	if err != nil {
		return nil, err
	}
	var templates []TaskTemplate
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		var template TaskTemplate
		err = json.NewDecoder(file).Decode(&template)
		if err != nil {
			return nil, err
		}
		if q != nil {
			panic("not implemented")
		}
		templates = append(templates, template)
	}
	return templates, nil
}

func (c Client) GenerateTasks() error {
	log.Infof("Generating tasks from task templates in %s", c.Config.WorkdirConfig.TaskTemplateDir)
	paths, err := findFiles(c.Config.WorkdirConfig.TaskTemplateDir, "*.json")
	if err != nil {
		return err
	}
	total := 0
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		var template TaskTemplate
		err = json.NewDecoder(file).Decode(&template)
		if err != nil {
			return err
		}
		task, err := template.GetTask(nil)
		if err != nil {
			return err
		}
		path := filepath.Join(c.Config.WorkdirConfig.TaskDir, fmt.Sprintf("%s.json", task.Id))
		file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		blob, err := json.MarshalIndent(task, "", "  ")
		if err != nil {
			return err
		}
		_, err = file.Write(blob)
		if err != nil {
			return err
		}
		total++
	}
	log.Infof("Wrote %d tasks to %s", total, c.Config.WorkdirConfig.TaskDir)
	return nil
}

func (c Client) GenerateTaskTemplates(atomicsPath string) error {
	templates, err := c.generateTaskTemplatesFromAtomicRedTeam(atomicsPath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(c.Config.WorkdirConfig.TaskTemplateDir, 0755)
	if err != nil {
		return err
	}
	for _, template := range templates {
		path := filepath.Join(c.Config.WorkdirConfig.TaskTemplateDir, fmt.Sprintf("%s.json", template.Id))
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		blob, err := json.MarshalIndent(template, "", "  ")
		if err != nil {
			return err
		}
		_, err = file.Write(blob)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Client) generateTaskTemplatesFromAtomicRedTeam(atomicsPath string) ([]TaskTemplate, error) {
	paths, err := findFiles(atomicsPath, "T*/T*.yaml")
	if err != nil {
		return nil, err
	}
	var result []TaskTemplate
	for _, path := range paths {
		bundle, err := ReadAtomicRedTeamYAMLFile(path)
		if err != nil {
			log.Warnf("Failed to parse bundle: %s - %s", path, err)
			continue
		}
		templates, err := bundle.GetTaskTemplates()
		if err != nil {
			log.Warnf("Failed to get task templates: %s - %s", path, err)
			continue
		}
		result = append(result, templates...)
	}
	return result, nil
}
