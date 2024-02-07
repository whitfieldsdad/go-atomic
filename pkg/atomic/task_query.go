package atomic

type TaskQuery struct {
	TaskIds           []string `json:"ids"`
	TaskTemplateIds   []string `json:"template_ids"`
	Tags              []string `json:"tags"`
	Platforms         []string `json:"platforms"`
	ElevationRequired *bool    `json:"elevation_required"`
}

func (q TaskQuery) MatchesTask(t Task) bool {
	if len(q.TaskIds) > 0 && !AnyStringMatchesAnyCaseInsensitivePattern([]string{t.Id}, q.TaskIds) {
		return false
	}
	if len(q.TaskTemplateIds) > 0 && !AnyStringMatchesAnyCaseInsensitivePattern([]string{t.TaskTemplateId}, q.TaskTemplateIds) {
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
