package configdata

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

// Validates the config against a schema (specific to config type)
func isValidInput(schemaJSONString, inputJSONString string) bool {
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", strings.NewReader(schemaJSONString)); err != nil {
		panic(err)
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		panic(err)
	}

	var doc interface{}
	if err := json.Unmarshal([]byte(inputJSONString), &doc); err != nil {
		return false
	}

	if err := schema.Validate(doc); err != nil {
		fmt.Println("Input not matches with the Schema:", err)
		return false
	}

	return true
}

// Validates two schemas properties is the same and ignore the order
func equalSchemas(a, b string) bool {
	var ma, mb map[string]interface{}

	if err := json.Unmarshal([]byte(a), &ma); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(b), &mb); err != nil {
		return false
	}

	return reflect.DeepEqual(ma, mb)
}
