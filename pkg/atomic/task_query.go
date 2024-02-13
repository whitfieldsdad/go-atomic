package atomic

type TaskTemplateQuery struct {
	TaskTemplateIds    []string `json:"template_ids"`
	AttackTechniqueIds []string `json:"attack_technique_ids"`
	Tags               []string `json:"tags"`
	Platforms          []string `json:"platforms"`
	ElevationRequired  *bool    `json:"elevation_required"`
}

func (q TaskTemplateQuery) Matches(t TaskTemplate) bool {
	if len(q.TaskTemplateIds) > 0 && !AnyStringMatchesAnyCaseInsensitivePattern([]string{t.Id}, q.TaskTemplateIds) {
		return false
	}
	if len(q.AttackTechniqueIds) > 0 && !AnyStringMatchesAnyCaseInsensitivePattern(t.GetAttackTechniqueIds(), q.AttackTechniqueIds) {
		return false
	}
	if len(q.Tags) > 0 && !AnyStringMatchesAnyCaseInsensitivePattern(t.GetTags(), q.Tags) {
		return false
	}
	if len(q.Platforms) > 0 && !AnyStringMatchesAnyCaseInsensitivePattern(t.Platforms, q.Platforms) {
		return false
	}
	if q.ElevationRequired != nil {
		elevationRequired := *q.ElevationRequired
		if elevationRequired != t.IsElevationRequired() {
			return false
		}
	}
	return true
}
