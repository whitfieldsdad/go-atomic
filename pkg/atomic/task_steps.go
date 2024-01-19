package atomic

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	StepTypeExec = "exec"
)

type Step struct {
	Id       string      `json:"id"`
	StepType string      `json:"step_type"`
	Data     interface{} `json:"data"`
}

func (s Step) Exec(ctx context.Context) StepResult {
	var (
		result interface{}
		err    error
	)
	startTime := time.Now()

	d := s.Data
	switch d.(type) {
	case ExecStep:
		o := d.(ExecStep)
		result, err = o.Exec(ctx)
	default:
		panic(fmt.Sprintf("Unresolved type: %s", s.StepType))
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	return StepResult{
		Id:        uuid.NewString(),
		StepId:    s.Id,
		StepType:  s.Id,
		StartTime: startTime,
		Duration:  duration.Seconds(),
		EndTime:   endTime,
		Data:      result,
		Error:     err,
	}
}

type StepResult struct {
	Id        string      `json:"id"`
	StepId    string      `json:"step_id"`
	StepType  string      `json:"step_type"`
	StartTime time.Time   `json:"start_time"`
	EndTime   time.Time   `json:"end_time"`
	Duration  float64     `json:"duration"`
	Data      interface{} `json:"data,omitempty"`
	Error     error       `json:"error,omitempty"`
}

type ExecStep struct {
	CommandTemplate string                 `json:"command_template"`
	CommandType     string                 `json:"command_type"`
	Args            map[string]interface{} `json:"args"`
}

func (s ExecStep) Exec(ctx context.Context) (*ExecStepResult, error) {
	c := NewShellCommandWithArgs(s.CommandTemplate, s.CommandType, s.Args)
	subprocess, err := c.Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &ExecStepResult{
		Subprocess: *subprocess,
	}, nil
}

type ExecStepResult struct {
	Subprocess Process `json:"subprocess,omitempty"`
}

func NewExecStep(commandTemplate, commandType string) (*Step, error) {
	return NewExecStepWithArgs(commandTemplate, commandType, nil)
}

func NewExecStepWithArgs(commandTemplate, commandType string, args map[string]interface{}) (*Step, error) {
	s := &Step{
		Id:       uuid.NewString(),
		StepType: StepTypeExec,
		Data: ExecStep{
			CommandTemplate: commandTemplate,
			CommandType:     commandType,
			Args:            args,
		},
	}
	return s, nil
}
