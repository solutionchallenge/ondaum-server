package utils

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type MarshalingMethod string

var (
	MarshalJSON = MarshalingMethod("MarshalJSON")
	MarshalYAML = MarshalingMethod("MarshalYAML")
)

func MustMarshal(v any, method MarshalingMethod) []byte {
	switch method {
	case MarshalJSON:
		json, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		return json
	case MarshalYAML:
		yaml, err := yaml.Marshal(v)
		if err != nil {
			panic(err)
		}
		return yaml
	default:
		panic(fmt.Sprintf("unknown marshaling method: %s", method))
	}
}
