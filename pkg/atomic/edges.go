package atomic

// Declare a string enum of edge types
type EdgeType string

const (
	EdgeTypeOnSuccess EdgeType = "on-success"
	EdgeTypeOnFailure EdgeType = "on-failure"
)

type Edge struct {
	Source   string   `json:"source"`
	Relation EdgeType `json:"relation"`
	Target   string   `json:"target"`
}
