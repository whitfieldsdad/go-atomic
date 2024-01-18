package atomic

import "time"

var (
	Created   = "created"
	Queued    = "queued"
	Running   = "running"
	Success   = "success"
	Failed    = "failed"
	Cancelled = "cancelled"
	Killed    = "killed"
)

type TaskInvocationStatus struct {
	Id               string    `json:"id"`
	TaskId           string    `json:"task_id"`
	TaskInvocationId string    `json:"task_invocation_id"`
	Time             time.Time `json:"time"`
	Status           string    `json:"status"`
}

func (s TaskInvocationStatus) IsCreated() bool {
	return s.Status == Created
}

func (s TaskInvocationStatus) IsQueued() bool {
	return s.Status == Queued
}

func (s TaskInvocationStatus) IsRunning() bool {
	return s.Status == Running
}

func (s TaskInvocationStatus) IsFailed() bool {
	return s.Status == Failed
}

func (s TaskInvocationStatus) IsSuccess() bool {
	return s.IsSuccess()
}

func (s TaskInvocationStatus) IsCancelled() bool {
	return s.Status == Cancelled
}

func (s TaskInvocationStatus) IsKilled() bool {
	return s.Status == Killed
}

func (s TaskInvocationStatus) IsDone() bool {
	return s.IsSuccess() || s.IsFailed() || s.IsCancelled() || s.IsKilled()
}
