package atomic

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

type FlowType string

const (
	FlowTypeOnSuccess FlowType = "on-success"
	FlowTypeOnFailure FlowType = "on-failure"
)

func (f FlowType) String() string {
	return string(f)
}

type Flow struct {
	S string   `json:"s"`
	P FlowType `json:"p"`
	O string   `json:"o"`
}

func NewFlow(s string, p FlowType, o string) Flow {
	return Flow{
		S: s,
		P: p,
		O: o,
	}
}

type TaskTemplate struct {
	Id             string          `json:"id"`
	Name           string          `json:"name,omitempty"`
	Description    string          `json:"description,omitempty"`
	Platforms      []string        `json:"platforms,omitempty"`
	Steps          []Step          `json:"steps"`
	Flows          []Flow          `json:"flows,omitempty"`
	Tags           []string        `json:"tags,omitempty"`
	InputArguments []InputArgument `json:"input_arguments,omitempty"`
}

func (t TaskTemplate) GetAttackTechniqueIds() []string {
	var result []string
	for _, tag := range t.Tags {
		if !slices.Contains(result, tag) && isAttackTechniqueId(tag) {
			result = append(result, tag)
		}
	}
	for _, step := range t.Steps {
		for _, tag := range step.Tags {
			if !slices.Contains(result, tag) && isAttackTechniqueId(tag) {
				result = append(result, tag)
			}
		}
	}
	return result
}

func (t TaskTemplate) IsElevationRequired() bool {
	for _, s := range t.Steps {
		if s.ElevationRequired {
			return true
		}
	}
	return false
}

func (t TaskTemplate) GetTags() []string {
	tags := t.Tags
	for _, step := range t.Steps {
		for _, tag := range step.Tags {
			if !slices.Contains(tags, tag) {
				tags = append(tags, tag)
			}
		}
	}
	return tags
}

func (t TaskTemplate) GetArgumentMap(args map[string]interface{}) map[string]interface{} {
	m := t.getDefaultArgumentMap()
	for k, v := range args {
		m[k] = v
	}
	return m
}

func (t TaskTemplate) getDefaultArgumentMap() map[string]interface{} {
	m := make(map[string]interface{})
	for _, arg := range t.InputArguments {
		m[arg.Name] = arg.Value
	}
	return m
}

func (t TaskTemplate) GetTask(args map[string]interface{}) (*Task, error) {
	m := t.GetArgumentMap(args)
	id, err := NewUUID5FromMap(t.Id, m)
	if err != nil {
		return nil, err
	}
	steps := make([]Step, len(t.Steps))
	for i, s := range t.Steps {
		if s.Type == StepTypeExecuteCommand {
			d := s.Data.(ExecuteCommandStep)
			d.Command = InterpolateCommandArgs(d.Command, m)
			steps[i] = Step{
				Id:   s.Id,
				Type: s.Type,
				Data: d,
				Tags: s.Tags,
			}
		}
	}
	return &Task{
		Id:          id,
		TemplateId:  t.Id,
		Name:        t.Name,
		Description: t.Description,
		Platforms:   t.Platforms,
		Steps:       t.Steps,
	}, nil
}

type Task struct {
	Id          string   `json:"id"`
	TemplateId  string   `json:"template_id"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Platforms   []string `json:"platforms,omitempty"`
	Steps       []Step   `json:"steps"`
	Flows       []Flow   `json:"flows,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func (t Task) IsElevationRequired() bool {
	for _, s := range t.Steps {
		if s.ElevationRequired {
			return true
		}
	}
	return false
}

func (t Task) IsRunnable() bool {
	if t.IsElevationRequired() && !IsElevated() {
		return false
	}
	if len(t.Platforms) > 0 && !slices.Contains(t.Platforms, runtime.GOOS) {
		return false
	}
	return true
}

func (t Task) GetTags() []string {
	tags := t.Tags
	for _, step := range t.Steps {
		for _, tag := range step.Tags {
			if !slices.Contains(tags, tag) {
				tags = append(tags, tag)
			}
		}
	}
	return tags
}

func (t Task) GetAttackTechniqueIds() []string {
	ids := make([]string, 0)
	for _, step := range t.Steps {
		for _, id := range step.GetAttackTechniqueIds() {
			if !slices.Contains(ids, id) {
				ids = append(ids, id)
			}
		}
	}
	return ids
}

func (t Task) Run(ctx context.Context) (*ExecutedTask, error) {
	if !t.IsRunnable() {
		return nil, errors.New("task is not runnable")
	}
	startTime := time.Now()
	var stepResults []StepResult
	for _, step := range t.Steps {
		stepResult := step.Run(ctx)
		stepResults = append(stepResults, stepResult)
	}
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	user, err := GetCurrentUser()
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup current user")
	}
	process, err := GetProcess(os.Getpid())
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup current process")
	}
	return &ExecutedTask{
		Id:        uuid.NewString(),
		HostId:    GetHostId(),
		AgentId:   GetAgentId(),
		UserId:    user.Id,
		TaskId:    t.Id,
		Steps:     stepResults,
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  duration.Seconds(),
		Process:   *process,
	}, nil
}

type ExecutedTask struct {
	Id        string       `json:"id"`
	TaskId    string       `json:"task_id"`
	AgentId   string       `json:"agent_id"`
	HostId    string       `json:"host_id"`
	UserId    string       `json:"user_id"`
	Steps     []StepResult `json:"steps"`
	StartTime time.Time    `json:"start_time"`
	EndTime   time.Time    `json:"end_time"`
	Duration  float64      `json:"duration"`
	Process   Process      `json:"process"`
}

func (e ExecutedTask) OK() bool {
	for _, step := range e.Steps {
		if !step.OK() {
			return false
		}
	}
	return true
}
