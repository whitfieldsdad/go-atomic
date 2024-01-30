package atomic

import (
	"regexp"
)

var (
	AttackTechniqueIdRegex = regexp.MustCompile("[Tt]\\d+(?:\\.\\d+)?")
	UUIDRegex              = regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
)

type TaskQuery struct {
	TaskIds            []string `json:"ids"`
	TaskTemplateIds    []string `json:"template_ids"`
	AttackTacticIds    []string `json:"attack_tactic_ids"`
	AttackTechniqueIds []string `json:"attack_technique_ids"`
	Tags               []string `json:"tags"`
	Platforms          []string `json:"platforms"`
	ElevationRequired  *bool    `json:"elevation_required"`
}

func (q TaskQuery) Matches(t Task) bool {
	if len(q.TaskIds) > 0 && !AnyStringMatchesAnyCaseInsensitivePattern([]string{t.Id}, q.TaskIds) {
		return false
	}
	if len(q.TaskTemplateIds) > 0 && !AnyStringMatchesAnyCaseInsensitivePattern([]string{t.TemplateId}, q.TaskTemplateIds) {
		return false
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
