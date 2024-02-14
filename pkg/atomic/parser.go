package atomic

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func parseTaskTemplate(b []byte) (*TaskTemplate, error) {
	var template TaskTemplate
	err := json.Unmarshal(b, &template)
	if err != nil {
		return nil, err
	}
	for i, s := range template.Steps {
		switch s.Type {
		case "execute-command":
			var d ExecuteCommandStep
			err := mapstructure.Decode(s.Data, &d)
			if err != nil {
				return nil, err
			}
			template.Steps[i].Data = d
		case "list-processes":
			var d ListProcessesStep
			err := mapstructure.Decode(s.Data, &d)
			if err != nil {
				return nil, err
			}
			template.Steps[i].Data = d
		case "list-users":
			var d ListUsersStep
			err := mapstructure.Decode(s.Data, &d)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("unknown step type")
		}
	}
	return &template, nil
}
