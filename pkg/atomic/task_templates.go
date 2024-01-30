package atomic

import (
	"context"
	"errors"
)

type TaskTemplate struct {
	Id                 string          `json:"id"`
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	Platforms          []string        `json:"platforms"`
	ElevationRequired  bool            `json:"elevation_required"`
	InputArguments     []InputArgument `json:"input_arguments"`
	AttackTechniqueIds []string        `json:"attack_technique_ids"`
	Tags               []string        `json:"tags"`
	Steps              []Step          `json:"steps"`
	Edges              []Edge          `json:"edges,omitempty"`
}

// TODO: pass input arguments
func (t TaskTemplate) GetTask() (*Task, error) {
	m := make(map[string]interface{})
	for _, arg := range t.InputArguments {
		m[arg.Name] = arg.Value
	}
	id, err := NewUUID5FromMap(t.Id, m)
	if err != nil {
		return nil, err
	}
	return &Task{
		Id:                 id,
		TemplateId:         t.Id,
		Name:               t.Name,
		Description:        t.Description,
		Platforms:          t.Platforms,
		ElevationRequired:  t.ElevationRequired,
		Steps:              t.Steps,
		Edges:              t.Edges,
		Tags:               t.Tags,
		AttackTechniqueIds: t.AttackTechniqueIds,
	}, nil
}

// TODO: pass input arguments
func (t *TaskTemplate) Exec(ctx context.Context) (*ExecutedTask, error) {
	task, err := t.GetTask()
	if err != nil {
		return nil, errors.New("failed to execute task from task template")
	}
	executedTask := task.Exec(ctx)
	return &executedTask, nil
}
