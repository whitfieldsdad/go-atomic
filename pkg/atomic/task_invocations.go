package atomic

import "time"

type TaskInvocation struct {
	Id     string    `json:"id"`
	Time   time.Time `json:"time"`
	TaskId string    `json:"task_id"`
	HostId string    `json:"host_id"`
	UserId string    `json:"user_id"`
}
