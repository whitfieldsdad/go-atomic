package atomic

import "time"

type TaskInvocation struct {
	Id     string    `json:"id"`
	TaskId string    `json:"task_id"`
	Time   time.Time `json:"time"`
}
