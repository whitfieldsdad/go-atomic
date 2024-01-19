package atomic

import (
	"os"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type AtomicRedTeamArgSpec struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
}

type AtomicRedTeamDependency struct {
	Description      string `yaml:"description"`
	PrereqCommand    string `yaml:"prereq_command"`
	GetPrereqCommand string `yaml:"get_prereq_command"`
	ExecutorName     string `yaml:"-"`
}

type AtomicRedTeamExecutor struct {
	Name              string `yaml:"name"`
	Command           string `yaml:"command"`
	CleanupCommand    string `yaml:"cleanup_command"`
	ElevationRequired bool   `yaml:"elevation_required"`
}

type AtomicRedTeamTest struct {
	Id                     string                          `yaml:"auto_generated_guid"`
	Name                   string                          `yaml:"name"`
	Executor               AtomicRedTeamExecutor           `yaml:"executor"`
	Description            string                          `yaml:"description"`
	Platforms              []string                        `yaml:"supported_platforms"`
	InputArguments         map[string]AtomicRedTeamArgSpec `yaml:"input_arguments,omitempty"`
	Dependencies           []AtomicRedTeamDependency       `yaml:"dependencies,omitempty"`
	DependencyExecutorName string                          `yaml:"dependency_exector_name,omitempty"`
	AttackTechniqueId      string                          `yaml:"-"`
}

func (t AtomicRedTeamTest) GetTask() (*Task, error) {
	return t.GetTaskWithArgs(nil)
}

func (t AtomicRedTeamTest) GetArgs() map[string]interface{} {
	m := make(map[string]interface{})
	for k, argspec := range t.InputArguments {
		v := argspec.Default
		m[k] = v
	}
	return m
}

func (t AtomicRedTeamTest) GetTaskWithArgs(args map[string]interface{}) (*Task, error) {
	if args == nil {
		args = t.GetArgs()
	} else {
		args = MergeInputArgs(t.GetArgs(), args)
	}

	var steps []Step
	testStep, err := NewExecStepWithArgs(t.Executor.Command, t.Executor.Name, args)
	if err != nil {
		return nil, err
	}
	steps = append(steps, *testStep)

	if t.Executor.CleanupCommand != "" {
		cleanupStep, err := NewExecStepWithArgs(t.Executor.CleanupCommand, t.Executor.Name, args)
		if err != nil {
			return nil, err
		}
		steps = append(steps, *cleanupStep)
	}

	task := &Task{
		Id:                 uuid.NewString(),
		Aliases:            []string{t.Id},
		AttackTechniqueIds: []string{t.AttackTechniqueId},
		Name:               t.Name,
		Description:        t.Description,
		Steps:              steps,
		Platforms:          t.Platforms,
		ElevationRequired:  t.Executor.ElevationRequired,
	}
	return task, nil
}

type AtomicRedTeamTestBundle struct {
	AttackTechniqueId string              `yaml:"attack_technique"`
	DisplayName       string              `yaml:"display_name"`
	Tests             []AtomicRedTeamTest `yaml:"atomic_tests"`
}

func (b AtomicRedTeamTestBundle) GetTasks() ([]Task, error) {
	var tasks []Task
	for _, test := range b.Tests {
		task, err := test.GetTask()
		if err != nil {
			continue
		}
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func ReadAtomicRedTeamYAMLFile(path string) (*AtomicRedTeamTestBundle, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var bundle AtomicRedTeamTestBundle
	err = yaml.Unmarshal(b, &bundle)
	if err != nil {
		return nil, err
	}
	var tests []AtomicRedTeamTest
	for _, test := range bundle.Tests {
		if test.Executor.Name == "manual" {
			continue
		}
		test.AttackTechniqueId = bundle.AttackTechniqueId
		for _, dependency := range test.Dependencies {
			dependency.ExecutorName = test.DependencyExecutorName
		}
		for i, platform := range test.Platforms {
			if platform == "macos" {
				test.Platforms[i] = "darwin"
			}
		}
		tests = append(tests, test)
	}
	bundle.Tests = tests
	return &bundle, nil
}
