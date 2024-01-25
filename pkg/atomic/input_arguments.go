package atomic

import (
	"fmt"
	"strings"
)

type ArgumentType string

const (
	ArgumentTypePath    ArgumentType = "path"
	ArgumentTypeString  ArgumentType = "string"
	ArgumentTypeInteger ArgumentType = "integer"
	ArgumentTypeFloat   ArgumentType = "float"
)

func (t ArgumentType) toString() string {
	return fmt.Sprintf("%s", t)
}

type InputArgument struct {
	Name        string       `json:"name"`
	Type        ArgumentType `json:"type"`
	Description string       `json:"description"`
	Value       interface{}  `json:"value"`
}

func (v InputArgument) IsPath() bool {
	return v.isType(v.Type, ArgumentTypePath)
}

func (v InputArgument) IsString() bool {
	return v.isType(v.Type, ArgumentTypePath)
}

func (v InputArgument) IsInteger() bool {
	return v.isType(v.Type, ArgumentTypeInteger)
}

func (v InputArgument) IsFloat() bool {
	return v.isType(v.Type, ArgumentTypeFloat)
}

func (v InputArgument) isType(a, b ArgumentType) bool {
	x := strings.ToLower(a.toString())
	y := strings.ToLower(b.toString())
	return x == y
}
