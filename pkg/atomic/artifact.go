package atomic

type ArtifactType string

const (
	ArtifactTypeFile    ArtifactType = "file"
	ArtifactTypeProcess ArtifactType = "process"
)

func (a ArtifactType) Plural() PluralArtifactType {
	return ArtifactTypeToPlural[a]
}

func (a ArtifactType) String() string {
	return string(a)
}

type PluralArtifactType string

const (
	PluralArtifactTypeFiles     PluralArtifactType = "files"
	PluralArtifactTypeProcesses PluralArtifactType = "processes"
)

var ArtifactTypeToPlural = map[ArtifactType]PluralArtifactType{
	ArtifactTypeFile:    PluralArtifactTypeFiles,
	ArtifactTypeProcess: PluralArtifactTypeProcesses,
}

func (a PluralArtifactType) String() string {
	return string(a)
}

type Artifact struct {
	Id string `json:"id"`
}

func NewArtifact() Artifact {
	return Artifact{
		Id: NewUUID4(),
	}
}
