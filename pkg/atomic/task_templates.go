package atomic

type TaskTemplate struct {
	TaskMetadata   `json:",inline"`
	InputArguments []InputArgument `json:"input_arguments" mapstructure:"input_arguments"`
}

func (t TaskTemplate) GetArgumentMap(args map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for _, a := range t.InputArguments {
		m[a.Name] = a.Value
	}
	for k, v := range args {
		m[k] = v
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
	task := &Task{
		TaskMetadata:   t.TaskMetadata,
		TaskTemplateId: t.Id,
	}
	task.Id = id
	return task, nil
}
