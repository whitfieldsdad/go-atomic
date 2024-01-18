package atomic

import "time"

type TaskInvocationResult struct {
	Id               string    `json:"id"`
	Time             time.Time `json:"time"`
	TaskId           string    `json:"task_id"`
	TaskInvocationId string    `json:"task_invocation_id"`
	ProcessIds       []string  `json:"process_ids`
}