package template

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
)

func GenerateSchema() string {
	schema := jsonschema.Reflect(&Template{})
	schemaBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return ""
	}
	return string(schemaBytes)
}
