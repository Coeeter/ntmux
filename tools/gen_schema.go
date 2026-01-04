package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/coeeter/ntmux/internal/template"
	"github.com/invopop/jsonschema"
)

var (
	schemaFile = flag.String("schema", "schema.json", "Output file for the JSON schema")
)

func main() {
	flag.Parse()
	schema := jsonschema.Reflect(&template.Template{})
	schemaBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(*schemaFile, schemaBytes, 0644); err != nil {
		panic(err)
	}
}
