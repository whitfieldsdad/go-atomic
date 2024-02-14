package atomic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type StepType string

const (
	StepTypeExecuteCommand StepType = "execute-command"
	StepTypeListProcesses  StepType = "list-processes"
	StepTypeListUsers      StepType = "list-users"
)

type Step struct {
	Id                string      `json:"id" mapstructure:"id"`
	Name              string      `json:"name,omitempty" mapstructure:"name"`
	Description       string      `json:"description,omitempty" mapstructure:"description"`
	Type              StepType    `json:"type" mapstructure:"type"`
	Data              interface{} `json:"data" mapstructure:"data"`
	ElevationRequired bool        `json:"elevation_required" mapstructure:"elevation_required"`
	Tags              []string    `json:"tags,omitempty" mapstructure:"tags"`
}

func (s Step) GetAttackTechniqueIds() []string {
	ids := make([]string, 0)
	for _, tag := range s.Tags {
		if !slices.Contains(ids, tag) && isAttackTechniqueId(tag) {
			ids = append(ids, tag)
		}
	}
	return ids
}

func (s Step) Run(ctx context.Context) StepResult {
	var (
		result interface{}
		err    error
	)
	startTime := time.Now()
	d := s.Data
	switch d.(type) {
	case ExecuteCommandStep:
		o := d.(ExecuteCommandStep)
		result, err = o.Run(ctx)
	case ListProcessesStep:
		o := d.(ListProcessesStep)
		result, err = o.Run(ctx)
	case ListUsersStep:
		o := d.(ListUsersStep)
		result, err = o.Run(ctx)
	default:
		panic(fmt.Sprintf("Unresolved type: %s", s.Type))
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	if err != nil {
		log.Warnf("Failed to execute step: %s", err)
	}

	return StepResult{
		Id:        uuid.NewString(),
		StepId:    s.Id,
		StepType:  s.Type,
		StartTime: startTime,
		Duration:  duration.Seconds(),
		EndTime:   endTime,
		Data:      result,
		Error:     err,
	}
}

type StepResult struct {
	Id        string      `json:"id" mapstructure:"id"`
	StepId    string      `json:"step_id" mapstructure:"step_id"`
	StepType  StepType    `json:"step_type" mapstructure:"step_type"`
	StartTime time.Time   `json:"start_time" mapstructure:"start_time"`
	EndTime   time.Time   `json:"end_time" mapstructure:"end_time"`
	Duration  float64     `json:"duration" mapstructure:"duration"`
	Data      interface{} `json:"data,omitempty" mapstructure:"data"`
	Error     error       `json:"error,omitempty" mapstructure:"error"`
}

func (s StepResult) OK() bool {
	if s.Error != nil {
		return false
	}
	d := s.Data
	switch d.(type) {
	case ExecuteCommandStepResult:
		o := d.(ExecuteCommandStepResult)
		return o.Process.ExitCode != nil && *o.Process.ExitCode == 0
	case ListProcessesStepResult:
		o := d.(ListProcessesStepResult)
		return len(o.Processes) > 0
	case ListUsersStepResult:
		o := d.(ListUsersStepResult)
		return len(o.Users) > 0
	default:
		return true
	}
}

type ListProcessesStep struct{}

func NewListProcessesStep() (*Step, error) {
	s := &Step{
		Id:   uuid.NewString(),
		Type: StepTypeListProcesses,
		Data: ListProcessesStep{},
		Tags: []string{"T1057"},
	}
	return s, nil
}

func (s ListProcessesStep) Run(ctx context.Context) (*ListProcessesStepResult, error) {
	processes, err := ListProcesses()
	if err != nil {
		return nil, err
	}
	return &ListProcessesStepResult{
		Processes: processes,
	}, nil
}

type ListProcessesStepResult struct {
	Processes []Process `json:"processes"`
}

type ExecuteCommandStep struct {
	Command     string `json:"command" mapstructure:"command"`
	CommandType string `json:"command_type" mapstructure:"command_type"`
}

func NewExecuteCommandStep(command, commandType string) (*Step, error) {
	s := &Step{
		Id:   uuid.NewString(),
		Type: StepTypeExecuteCommand,
		Data: ExecuteCommandStep{
			Command:     strings.TrimSpace(command),
			CommandType: strings.TrimSpace(commandType),
		},
	}
	return s, nil
}

func (s ExecuteCommandStep) Run(ctx context.Context) (*ExecuteCommandStepResult, error) {
	if s.Command == "" {
		return nil, fmt.Errorf("command is required")
	}
	if s.CommandType == "" {
		return nil, fmt.Errorf("command type is required")
	}
	c := NewShellCommand(s.Command, s.CommandType)
	process, err := c.Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &ExecuteCommandStepResult{
		Command:     c.Command,
		CommandType: c.CommandType,
		Process:     *process,
	}, nil
}

type ExecuteCommandStepResult struct {
	Command     string  `json:"command" mapstructure:"command"`
	CommandType string  `json:"command_type" mapstructure:"command_type"`
	Process     Process `json:"process" mapstructure:"process"`
}

type ListUsersStep struct{}

func NewListUsersStep() (*Step, error) {
	s := &Step{
		Id:   uuid.NewString(),
		Type: StepTypeListUsers,
		Data: ListUsersStep{},
		Tags: []string{},
	}
	return s, nil
}

func (s ListUsersStep) Run(ctx context.Context) (*ListUsersStepResult, error) {
	users, err := ListUsers()
	if err != nil {
		return nil, err
	}
	return &ListUsersStepResult{
		Users: users,
	}, nil
}

type ListUsersStepResult struct {
	Users []User `json:"users"`
}
