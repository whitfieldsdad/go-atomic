package atomic

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id                string   `json:"id"`
	Aliases           []string `json:"aliases,omitempty"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Steps             []Step   `json:"steps"`
	Platforms         []string `json:"platforms"`
	ElevationRequired bool     `json:"elevation_required"`
	Tags              []string `json:"tags"`
}

func NewTask(id string, steps []Step) *Task {
	return &Task{Id: id, Steps: steps}
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

	return ExecutedTask{
		Id:            uuid.NewString(),
		TaskId:        t.Id,
		ExecutedSteps: stepResults,
		StartTime:     startTime,
		EndTime:       endTime,
		Duration:      duration.Seconds(),
	}
}

type ExecutedTask struct {
	Id            string       `json:"id"`
	TaskId        string       `json:"task_id"`
	ExecutedSteps []StepResult `json:"steps"`
	StartTime     time.Time    `json:"start_time"`
	EndTime       time.Time    `json:"end_time"`
	Duration      float64      `json:"duration"`
}
