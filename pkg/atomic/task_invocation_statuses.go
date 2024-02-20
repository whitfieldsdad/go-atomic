package atomic

import (
	"time"

	"github.com/google/uuid"
)

type TaskInvocationStatusType string

const (
	TaskCreated TaskInvocationStatusType = "created"
	TaskQueued  TaskInvocationStatusType = "queued"
	TaskRunning TaskInvocationStatusType = "running"
	TaskFailed  TaskInvocationStatusType = "failed"
	TaskSuccess TaskInvocationStatusType = "success"
)

func (t TaskInvocationStatusType) String() string {
	return string(t)
}

type TaskInvocationStatus struct {
	Id               string                   `json:"id"`
	Time             time.Time                `json:"time"`
	StatusType       TaskInvocationStatusType `json:"status_type"`
	TaskId           string                   `json:"task_id"`
	TaskInvocationId string                   `json:"task_invocation_id"`
}

func (s TaskInvocationStatus) IsCreated() bool {
	return s.StatusType == TaskCreated
}

func (s TaskInvocationStatus) IsQueued() bool {
	return s.StatusType == TaskQueued
}

func (s TaskInvocationStatus) IsRunning() bool {
	return s.StatusType == TaskRunning
}

func (s TaskInvocationStatus) IsFailed() bool {
	return s.StatusType == TaskFailed
}

func (s TaskInvocationStatus) IsSuccess() bool {
	return s.StatusType == TaskSuccess
}

func (s TaskInvocationStatus) IsDone() bool {
	return s.IsSuccess() || s.IsFailed()
}

func NewTaskInvocationStatus(taskId, taskStepId, taskInvocationId, status string) TaskInvocationStatus {
	return TaskInvocationStatus{
		Id:               uuid.NewString(),
		Time:             time.Now(),
		StatusType:       TaskInvocationStatusType(status),
		TaskId:           taskId,
		TaskInvocationId: taskInvocationId,
	}
}
