package atomic

import (
	"time"

	"github.com/google/uuid"
)

var (
	Created   = "created"
	Queued    = "queued"
	Running   = "running"
	Success   = "success"
	Failed    = "failed"
	Cancelled = "cancelled"
	Killed    = "killed"
)

type TaskInvocationStatusInterface interface {
	GetSubtype() string
}

type baseTaskInvocationStatus struct {
	Id               string    `json:"id"`
	Time             time.Time `json:"time"`
	Status           string    `json:"status"`
	TaskId           string    `json:"task_id"`
	TaskInvocationId string    `json:"task_invocation_id"`
	HostId           string    `json:"host_id"`
	UserId           string    `json:"user_id"`
}

func (s baseTaskInvocationStatus) IsCreated() bool {
	return s.Status == Created
}

func (s baseTaskInvocationStatus) IsQueued() bool {
	return s.Status == Queued
}

func (s baseTaskInvocationStatus) IsRunning() bool {
	return s.Status == Running
}

func (s baseTaskInvocationStatus) IsFailed() bool {
	return s.Status == Failed
}

func (s baseTaskInvocationStatus) IsSuccess() bool {
	return s.IsSuccess()
}

func (s baseTaskInvocationStatus) IsCancelled() bool {
	return s.Status == Cancelled
}

func (s baseTaskInvocationStatus) IsKilled() bool {
	return s.Status == Killed
}

func (s baseTaskInvocationStatus) IsDone() bool {
	return s.IsSuccess() || s.IsFailed() || s.IsCancelled() || s.IsKilled()
}

type TaskInvocationStatus struct {
	baseTaskInvocationStatus
}

func (s TaskInvocationStatus) GetSubtype() string {
	return "task"
}

type TaskStepInvocationStatus struct {
	baseTaskInvocationStatus
	TaskStepId string `json:"task_step_id"`
}

func (s TaskStepInvocationStatus) GetSubtype() string {
	return "step"
}

func NewTaskInvocationStatus(taskId, taskStepId, taskInvocationId, hostId, userId, status string) TaskInvocationStatus {
	return TaskInvocationStatus{
		baseTaskInvocationStatus: baseTaskInvocationStatus{
			Id:               uuid.NewString(),
			Time:             time.Now(),
			Status:           status,
			TaskId:           taskId,
			TaskInvocationId: taskInvocationId,
			HostId:           hostId,
			UserId:           userId,
		},
	}
}

func NewTaskStepInvocationStatus(taskId, taskStepId, taskInvocationId, hostId, userId, status string) TaskStepInvocationStatus {
	return TaskStepInvocationStatus{
		baseTaskInvocationStatus: baseTaskInvocationStatus{
			Id:               uuid.NewString(),
			Time:             time.Now(),
			Status:           status,
			TaskId:           taskId,
			TaskInvocationId: taskInvocationId,
			HostId:           hostId,
			UserId:           userId,
		},
		TaskStepId: taskStepId,
	}
}
