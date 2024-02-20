package atomic

import (
	"runtime"

	"golang.org/x/exp/slices"
)

type FlowType string

const (
	FlowTypeOnSuccess FlowType = "on-success"
	FlowTypeOnFailure FlowType = "on-failure"
)

func (f FlowType) String() string {
	return string(f)
}

type Flow struct {
	S string   `json:"s"`
	P FlowType `json:"p"`
	O string   `json:"o"`
}

func NewFlow(s string, p FlowType, o string) Flow {
	return Flow{
		S: s,
		P: p,
		O: o,
	}
}

type TaskMetadata struct {
	Id          string   `json:"id" mapstructure:"id"`
	Name        string   `json:"name" mapstructure:"name"`
	Description string   `json:"description" mapstructure:"description"`
	Platforms   []string `json:"platforms" mapstructure:"platforms"`
	Steps       []Step   `json:"steps" mapstructure:"steps"`
	Flows       []Flow   `json:"flows" mapstructure:"flows"`
	Tags        []string `json:"tags" mapstructure:"tags"`
}

func (t TaskMetadata) GetTags() []string {
	tags := t.Tags
	for _, step := range t.Steps {
		for _, tag := range step.Tags {
			if !slices.Contains(tags, tag) {
				tags = append(tags, tag)
			}
		}
	}
	return tags
}

func (t TaskMetadata) IsElevationRequired() bool {
	for _, s := range t.Steps {
		if s.ElevationRequired {
			return true
		}
	}
	return false
}

func (t Task) IsRunnable() bool {
	if t.IsElevationRequired() && !IsElevated() {
		return false
	}
	if len(t.Platforms) > 0 && !slices.Contains(t.Platforms, runtime.GOOS) {
		return false
	}
	return true
}
