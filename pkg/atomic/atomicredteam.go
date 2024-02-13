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

	tags := []string{
		t.AttackTechniqueId,
		"atomic-red-team",
	}
	steps := make([]Step, 0)
	flows := make([]Flow, 0)

	// Create the test and cleanup steps.
	testStep, err := NewExecuteCommandStep(t.Executor.Command, t.Executor.Name)
	if err != nil {
		return nil, err
	}
	steps = append(steps, *testStep)

	if t.Executor.CleanupCommand != "" {
		cleanupStep, err := NewExecuteCommandStep(t.Executor.CleanupCommand, t.Executor.Name)
		if err != nil {
			return nil, err
		}
		flows = append(flows, NewFlow(testStep.Id, FlowTypeOnSuccess, cleanupStep.Id))
		flows = append(flows, NewFlow(testStep.Id, FlowTypeOnFailure, cleanupStep.Id))
		steps = append(steps, *cleanupStep)
	}

	// Check, resolve, and re-check any dependencies before running the test or cleanup commands.
	n := len(t.Dependencies)
	if n > 0 {
		for i, d := range t.Dependencies {
			a, err := NewExecuteCommandStep(d.PrereqCommand, d.ExecutorName)
			if err != nil {
				return nil, err
			}
			a.Name = fmt.Sprintf("Check dependency %d/%d", i+1, n)
			a.Description = d.Description

			b, err := NewExecuteCommandStep(d.GetPrereqCommand, d.ExecutorName)
			if err != nil {
				return nil, err
			}
			b.Name = fmt.Sprintf("Resolve dependency %d/%d", i+1, n)
			b.Description = d.Description

			c, err := NewExecuteCommandStep(d.PrereqCommand, d.ExecutorName)
			if err != nil {
				return nil, err
			}
			c.Name = fmt.Sprintf("Re-check dependency %d/%d", i+1, n)
			c.Description = d.Description

			steps = append(steps, *a, *b, *c)
			flows = append(flows, NewFlow(a.Id, FlowTypeOnSuccess, testStep.Id))
			flows = append(flows, NewFlow(a.Id, FlowTypeOnFailure, b.Id))
			flows = append(flows, NewFlow(b.Id, FlowTypeOnSuccess, c.Id))
			flows = append(flows, NewFlow(c.Id, FlowTypeOnFailure, testStep.Id))
		}
	}

	task := &TaskTemplate{
		Id:             t.Id,
		Name:           t.Name,
		Description:    t.Description,
		Platforms:      t.Platforms,
		Steps:          steps,
		Flows:          flows,
		Tags:           tags,
		InputArguments: args,
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

func readAtomicRedTeamYAMLFile(path string) (*AtomicRedTeamTestBundle, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseAtomicRedTeamYAML(b)
}

func parseAtomicRedTeamYAML(b []byte) (*AtomicRedTeamTestBundle, error) {
	var bundle AtomicRedTeamTestBundle
	err := yaml.Unmarshal(b, &bundle)
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
