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

type Task struct {
	Id                 string   `json:"id"`
	TemplateId         string   `json:"template_id"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	Platforms          []string `json:"platforms"`
	ElevationRequired  bool     `json:"elevation_required"`
	AttackTechniqueIds []string `json:"attack_technique_ids"`
	Tags               []string `json:"tags"`
	Steps              []Step   `json:"steps"`
	Edges              []Edge   `json:"edges,omitempty"`
}

func NewTask(id string, steps []Step) *Task {
	return &Task{Id: id, Steps: steps}
}

func (t *Task) IsRunnable() bool {
	if t.ElevationRequired && !IsElevated() {
		return false
	}
	if len(t.Platforms) > 0 && !slices.Contains(t.Platforms, runtime.GOOS) {
		return false
	}
	return true
}

func (t *Task) Exec(ctx context.Context) ExecutedTask {
	startTime := time.Now()
	var stepResults []StepResult
	for _, step := range t.Steps {
		stepResult := step.Exec(ctx)
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
