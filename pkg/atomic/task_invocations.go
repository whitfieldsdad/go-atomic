package atomic

import "time"

type TaskInvocation struct {
	Id     string    `json:"id"`
	Time   time.Time `json:"time"`
	TaskId string    `json:"task_id"`
}

func NewTaskInvocation(taskId string) *TaskInvocation {
	return &TaskInvocation{
		Id:     NewUUID4(),
		Time:   time.Now(),
		TaskId: taskId,
	}
}
