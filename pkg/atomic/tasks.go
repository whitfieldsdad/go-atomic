package atomic

type Task struct {
	TaskMetadata   `json:",inline"`
	TaskTemplateId string `json:"task_template_id"`
}
