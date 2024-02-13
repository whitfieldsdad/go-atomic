package atomic

type TaskMetadata struct {
	Id                string   `json:"id"`
	Name              string   `json:"name,omitempty"`
	Description       string   `json:"description,omitempty"`
	Platforms         []string `json:"platforms,omitempty"`
	ElevationRequired bool     `json:"elevation_required"`
	Tags              []string `json:"tags,omitempty"`
	Steps             []Step   `json:"steps"`
	OnSuccess         []string `json:"on_success,omitempty"`
	OnFailure         []string `json:"on_failure,omitempty"`
}
