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
	TaskIds            []string `json:"ids"`
	AttackTechniqueIds []string `json:"attack_technique_ids"`
	Tags               []string `json:"tags"`
	Platforms          []string `json:"platforms"`
	ElevationRequired  *bool    `json:"elevation_required"`
}

func (q TaskQuery) Matches(t Task) bool {
	if len(q.TaskIds) > 0 {
		aliases := []string{t.Id}
		if len(t.Aliases) > 0 {
			aliases = append(aliases, t.Aliases...)
		}
		if !AnyStringMatchesAnyCaseInsensitivePattern(aliases, q.TaskIds) {
			return false
		}
	}
	if len(q.AttackTechniqueIds) > 0 && !AnyStringMatchesAnyCaseInsensitivePattern(t.AttackTechniqueIds, q.AttackTechniqueIds) {
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

func isPlatform(value string) bool {
	return slices.Contains(platforms, value)
}

func isAttackTechniqueId(value string) bool {
	return AttackTechniqueIdRegex.MatchString(value)
}

func isUUID(value string) bool {
	return UUIDRegex.MatchString(value)
}
