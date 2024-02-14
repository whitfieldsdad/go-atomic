package atomic

import "time"

type TaskInvocationResult struct {
	Id        string       `json:"id"`
	Time      time.Time    `json:"time"`
	TaskId    string       `json:"task_id"`
	StartTime time.Time    `json:"start_time"`
	EndTime   time.Time    `json:"end_time"`
	Duration  float64      `json:"duration"`
	Process   Process      `json:"process"`
	User      User         `json:"user"`
	Steps     []StepResult `json:"steps"`
}
