package atomic

import (
	"regexp"

	"golang.org/x/exp/slices"
)

var (
	AttackTechniqueIdRegex = regexp.MustCompile("[Tt]\\d+(?:\\.\\d+)?")
	UUIDRegex              = regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
)

var (
	platforms = []string{
		"windows",
		"linux",
		"darwin",
	}
)

type TaskQuery struct {
	TaskIds           []string `json:"ids"`
	Tags              []string `json:"tags"`
	Platforms         []string `json:"platforms"`
	ElevationRequired *bool    `json:"elevation_required"`
}

func (q TaskQuery) Matches(t Task) bool {
	if len(q.TaskIds) > 0 && !StringMatchesAnyCaseInsensitivePattern(t.Id, q.TaskIds) {
		return false
	}
	if len(q.Tags) > 0 && !AnyStringMatchesAnyCaseInsensitivePattern(t.Tags, q.Tags) {
		return false
	}
	if len(q.Platforms) > 0 && !AnyStringMatchesAnyCaseInsensitivePattern(t.Platforms, q.Platforms) {
		return false
	}
	if q.ElevationRequired != nil {
		elevationRequired := *q.ElevationRequired
		if elevationRequired != t.ElevationRequired {
			return false
		}
	}
	return true
}

func (q TaskQuery) Merge(o TaskQuery) {
	q.TaskIds = append(q.TaskIds, o.TaskIds...)
	q.Tags = append(q.Tags, o.Tags...)
	q.Platforms = append(q.Platforms, o.Platforms...)
}

func isPlatform(value string) bool {
	return slices.Contains(platforms, value)
}

func isAttackTechniqueId(value string) bool {
	return AttackTechniqueIdRegex.MatchString(value)
}

func isUUID(value string) bool {
	return UUIDRegex.MatchString(value)
}
