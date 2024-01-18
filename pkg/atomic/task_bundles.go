package atomic

import "context"

type TaskBundle struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Tasks       []Task `json:"tasks"`
}

func (b TaskBundle) Exec(ctx context.Context) []ExecutedTask {
	var results []ExecutedTask
	for _, t := range b.Tasks {
		result := t.Exec(ctx)
		results = append(results, result)
	}
	return results
}
