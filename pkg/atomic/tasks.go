package atomic

import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type TaskMetadata struct {
	Id                string   `json:"id"`
	Name              string   `json:"name,omitempty"`
	Description       string   `json:"description,omitempty"`
	Platforms         []string `json:"platforms,omitempty"`
	ElevationRequired bool     `json:"elevation_required"`
	Tags              []string `json:"tags,omitempty"`
	Steps             []Step   `json:"steps"`
	OnSuccess         []string `json:"on_success,omitempty"`
	OnFailure         []string `json:"on_failure,omitempty"`
}

type TaskTemplate struct {
	TaskMetadata   `json:",inline"`
	InputArguments []InputArgument `json:"input_arguments,omitempty"`
}

// GetArgumentMap returns a map of input arguments with default values overridden by any provided arguments.
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
	return &Task{
		TaskMetadata: TaskMetadata{
			Id:                id,
			Name:              t.Name,
			Description:       t.Description,
			Platforms:         t.Platforms,
			ElevationRequired: t.ElevationRequired,
			Tags:              t.Tags,
			Steps:             t.Steps,
			OnSuccess:         t.OnSuccess,
			OnFailure:         t.OnFailure,
		},
		TaskTemplateId: t.Id,
	}, nil
}

type Task struct {
	TaskMetadata   `json:",inline"`
	TaskTemplateId string `json:"task_template_id,omitempty"`
}

// IsRunnable checks if the task is runnable on the current platform and with the current permissions.
func (t Task) IsRunnable() bool {
	if t.ElevationRequired && !IsElevated() {
		return false
	}
	if len(t.Platforms) > 0 && !slices.Contains(t.Platforms, runtime.GOOS) {
		return false
	}
	return true
}

// GetTags returns a list of unique tags across all steps in the task.
func (t Task) GetTags() []string {
	tags := make([]string, 0)
	for _, step := range t.Steps {
		for _, tag := range step.Tags {
			if !slices.Contains(tags, tag) {
				tags = append(tags, tag)
			}
		}
	}
	return tags
}

// GetAttackTechniqueIds returns a list of unique ATT&CK technique IDs across all steps in the task.
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

func (t Task) Run(ctx context.Context) ExecutedTask {
	log.Infof("Running task %s", t.Id)

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
		log.Fatalf("Failed to lookup current user: %s", err)
	}
	process, err := GetProcess(os.Getpid())
	if err != nil {
		log.Fatalf("Failed to lookup current process: %s", err)
	}
	log.Infof("Task %s completed", t.Id)
	return ExecutedTask{
		Id:        uuid.NewString(),
		HostId:    GetHostId(),
		UserId:    user.Id,
		TaskId:    t.Id,
		Steps:     stepResults,
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  duration.Seconds(),
		Process:   *process,
	}
}

type ExecutedTask struct {
	Id        string       `json:"id"`
	TaskId    string       `json:"task_id"`
	HostId    string       `json:"host_id"`
	UserId    string       `json:"user_id"`
	Steps     []StepResult `json:"steps"`
	StartTime time.Time    `json:"start_time"`
	EndTime   time.Time    `json:"end_time"`
	Duration  float64      `json:"duration"`
	Process   Process      `json:"process"`
}

// OK checks if all steps in the task completed successfully.
func (e ExecutedTask) OK() bool {
	for _, step := range e.Steps {
		if !step.OK() {
			return false
		}
	}
	return true
}
