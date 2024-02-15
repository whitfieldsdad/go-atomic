package atomic

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
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

func (t AtomicRedTeamTest) ReferencesAtomicsFolder() bool {
	haystack := []string{
		t.Executor.Command,
		t.Executor.CleanupCommand,
	}
	for _, d := range t.Dependencies {
		haystack = append(haystack, d.PrereqCommand, d.GetPrereqCommand)
	}
	for _, arg := range t.InputArguments {
		haystack = append(haystack, arg.Default)
	}
	for _, v := range haystack {
		if strings.Contains(v, "PathToAtomicsFolder") {
			return true
		}
	}
	return false
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

	tags := make([]string, 0)
	if strings.Contains(t.AttackTechniqueId, ".") {
		parentTechniqueId := strings.Split(t.AttackTechniqueId, ".")[0]
		tags = append(tags, parentTechniqueId, t.AttackTechniqueId)
	} else {
		tags = append(tags, t.AttackTechniqueId)
	}

	steps := make([]Step, 0)
	flows := make([]Flow, 0)

	// Create the test and cleanup steps.
	testStep, err := NewExecuteCommandStep(t.Executor.Command, t.Executor.Name)
	if err != nil {
		return nil, err
	}
	testStep.ElevationRequired = t.Executor.ElevationRequired
	steps = append(steps, *testStep)

	if t.Executor.CleanupCommand != "" {
		cleanupStep, err := NewExecuteCommandStep(t.Executor.CleanupCommand, t.Executor.Name)
		if err != nil {
			return nil, err
		}
		cleanupStep.ElevationRequired = t.Executor.ElevationRequired
		flows = append(flows, NewFlow(testStep.Id, FlowTypeOnSuccess, cleanupStep.Id))
		flows = append(flows, NewFlow(testStep.Id, FlowTypeOnFailure, cleanupStep.Id))
		steps = append(steps, *cleanupStep)
	}

	// Check, resolve, and re-check any dependencies before running the test or cleanup commands.
	n := len(t.Dependencies)
	if n > 0 {
		dependencyExecutorName := t.DependencyExecutorName
		if dependencyExecutorName == "" {
			dependencyExecutorName = t.Executor.Name
		}
		for i, d := range t.Dependencies {
			a, err := NewExecuteCommandStep(d.PrereqCommand, dependencyExecutorName)
			if err != nil {
				return nil, err
			}
			a.Name = fmt.Sprintf("Check dependency %d/%d", i+1, n)
			a.Description = d.Description
			a.ElevationRequired = t.Executor.ElevationRequired

			b, err := NewExecuteCommandStep(d.GetPrereqCommand, dependencyExecutorName)
			if err != nil {
				return nil, err
			}
			b.Name = fmt.Sprintf("Resolve dependency %d/%d", i+1, n)
			b.Description = d.Description
			b.ElevationRequired = t.Executor.ElevationRequired

			c, err := NewExecuteCommandStep(d.PrereqCommand, dependencyExecutorName)
			if err != nil {
				return nil, err
			}
			c.Name = fmt.Sprintf("Re-check dependency %d/%d", i+1, n)
			c.Description = d.Description
			c.ElevationRequired = t.Executor.ElevationRequired

			steps = append(steps, *a, *b, *c)
			flows = append(flows, NewFlow(a.Id, FlowTypeOnSuccess, testStep.Id))
			flows = append(flows, NewFlow(a.Id, FlowTypeOnFailure, b.Id))
			flows = append(flows, NewFlow(b.Id, FlowTypeOnSuccess, c.Id))
			flows = append(flows, NewFlow(c.Id, FlowTypeOnFailure, testStep.Id))
		}
	}

	task := &TaskTemplate{
		TaskMetadata: TaskMetadata{
			Id:          t.Id,
			Name:        t.Name,
			Description: t.Description,
			Platforms:   t.Platforms,
			Steps:       steps,
			Flows:       flows,
			Tags: tags,
		},
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
		if test.ReferencesAtomicsFolder() {
			log.Debugf("Skipping test %s because it references the atomics folder (not supported)", test.Id)
			continue
		}
		test.AttackTechniqueId = bundle.AttackTechniqueId
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
