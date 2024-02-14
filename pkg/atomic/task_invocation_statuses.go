package atomic

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatusType string

const (
	TaskCreated TaskStatusType = "created"
	TaskQueued  TaskStatusType = "queued"
	TaskRunning TaskStatusType = "running"
	TaskFailed  TaskStatusType = "failed"
	TaskSuccess TaskStatusType = "success"
)

func (t TaskStatusType) String() string {
	return string(t)
}

type TaskInvocationStatusInterface interface {
	GetSubtype() string
}

type baseTaskInvocationStatus struct {
	Id               string    `json:"id"`
	Time             time.Time `json:"time"`
	Status           string    `json:"status"`
	TaskId           string    `json:"task_id"`
	TaskInvocationId string    `json:"task_invocation_id"`
}

func (s baseTaskInvocationStatus) IsCreated() bool {
	return s.Status == string(TaskCreated)
}

func (s baseTaskInvocationStatus) IsQueued() bool {
	return s.Status == string(TaskQueued)
}

func (s baseTaskInvocationStatus) IsRunning() bool {
	return s.Status == string(TaskRunning)
}

func (s baseTaskInvocationStatus) IsFailed() bool {
	return s.Status == string(TaskFailed)
}

func (s baseTaskInvocationStatus) IsSuccess() bool {
	return s.IsSuccess()
}

func (s baseTaskInvocationStatus) IsDone() bool {
	return s.IsSuccess() || s.IsFailed()
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

func NewTaskInvocationStatus(taskId, taskStepId, taskInvocationId, status string) TaskInvocationStatus {
	return TaskInvocationStatus{
		baseTaskInvocationStatus: baseTaskInvocationStatus{
			Id:               uuid.NewString(),
			Time:             time.Now(),
			Status:           status,
			TaskId:           taskId,
			TaskInvocationId: taskInvocationId,
		},
	}
}

func NewTaskStepInvocationStatus(taskId, taskStepId, taskInvocationId, status string) TaskStepInvocationStatus {
	return TaskStepInvocationStatus{
		baseTaskInvocationStatus: baseTaskInvocationStatus{
			Id:               uuid.NewString(),
			Time:             time.Now(),
			Status:           status,
			TaskId:           taskId,
			TaskInvocationId: taskInvocationId,
		},
		TaskStepId: taskStepId,
	}
}
