package atomic

import "time"

type Artifact struct {
	Id   string    `json:"id"`
	Time time.Time `json:"time"`
}

func NewArtifact() Artifact {
	return Artifact{
		Id:   NewUUID4(),
		Time: time.Now(),
	}
}
