package atomic

type TaskTemplate struct {
	Id                string          `json:"id"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	Platforms         []string        `json:"platforms"`
	ElevationRequired bool            `json:"elevation_required"`
	InputArguments    []InputArgument `json:"input_arguments"`
	Tags              []string        `json:"tags"`
	Steps             []Step          `json:"steps"`
	Edges             []string        `json:"edges"`
}
