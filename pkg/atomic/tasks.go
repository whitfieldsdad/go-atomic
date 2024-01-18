package atomic

import (
	"context"

	"github.com/google/uuid"
)

type Task struct {
	Id                string   `json:"id"`
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
	et := ExecutedTask{
		Id:     uuid.NewString(),
		TaskId: t.Id,
	}
	for _, s := range t.Steps {
		es := s.Exec(ctx)
		et.ExecutedSteps = append(et.ExecutedSteps, es)
	}
	return et
}

type ExecutedTask struct {
	Id            string       `json:"id"`
	TaskId        string       `json:"task_id"`
	ExecutedSteps []StepResult `json:"steps"`
}
