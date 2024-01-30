package atomic

import (
	"fmt"
	"os"

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

func (t AtomicRedTeamTest) GetArgs() map[string]interface{} {
	m := make(map[string]interface{})
	for k, argspec := range t.InputArguments {
		v := argspec.Default
		m[k] = v
	}
	return m
}

// TODO: add dependency check steps
// TODO: add dependency resolution steps
func (t AtomicRedTeamTest) GetTaskTemplate() (*TaskTemplate, error) {
	var args []InputArgument
	for k, argspec := range t.InputArguments {
		arg := InputArgument{
			Name:        k,
			Type:        ArgumentType(argspec.Type),
			Description: argspec.Description,
			Value:       argspec.Default,
		}
		args = append(args, arg)
	}

	var steps []Step
	var edges []Edge

	testStep, err := NewExecStep(t.Executor.Command, t.Executor.Name)
	if err != nil {
		return nil, err
	}
	steps = append(steps, *testStep)

	// Always execute cleanup commands.
	if t.Executor.CleanupCommand != "" {
		cleanupStep, err := NewExecStep(t.Executor.CleanupCommand, t.Executor.Name)
		if err != nil {
			return nil, err
		}
		steps = append(steps, *cleanupStep)

		edges = append(edges, Edge{
			Source:   testStep.Id,
			Relation: EdgeTypeOnSuccess,
			Target:   cleanupStep.Id,
		})
		edges = append(edges, Edge{
			Source:   testStep.Id,
			Relation: EdgeTypeOnFailure,
			Target:   cleanupStep.Id,
		})
	}

	task := &TaskTemplate{
		Id:                 t.Id,
		Name:               t.Name,
		DisplayName:        fmt.Sprintf("%s: %s", t.AttackTechniqueId, t.Name),
		Description:        t.Description,
		Platforms:          t.Platforms,
		ElevationRequired:  t.Executor.ElevationRequired,
		InputArguments:     args,
		Tags:               []string{t.AttackTechniqueId},
		AttackTechniqueIds: []string{t.AttackTechniqueId},
		Steps:              steps,
		Edges:              edges,
	}
	return task, nil
}

type AtomicRedTeamTestBundle struct {
	AttackTechniqueId string              `yaml:"attack_technique"`
	DisplayName       string              `yaml:"display_name"`
	Tests             []AtomicRedTeamTest `yaml:"atomic_tests"`
}

func (b AtomicRedTeamTestBundle) GetTaskTemplates() ([]TaskTemplate, error) {
	var tasks []TaskTemplate
	for _, test := range b.Tests {
		task, err := test.GetTaskTemplate()
		if err != nil {
			continue
		}
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func (b AtomicRedTeamTestBundle) GetTasks() ([]Task, error) {
	var tasks []Task
	for _, test := range b.Tests {
		template, err := test.GetTaskTemplate()
		if err != nil {
			continue
		}
		task, err := template.GetTask()
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
