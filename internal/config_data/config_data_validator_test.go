package configdata

import "testing"

func TestIsValidInput_Success(t *testing.T) {
	schemaJSON := `{
		"type": "object",
		"properties": {
			"enabled": { "type": "boolean" },
			"max_limit": { "type": "integer" }
		},
		"required": ["enabled", "max_limit"]
	}`

	inputJSON := `{
		"enabled": true,
		"max_limit": 100
	}`

	if !isValidInput(schemaJSON, inputJSON) {
		t.Fatal("expected input to match schema, but validation failed")
	}
}

func TestIsValidInput_MissingRequiredProperty(t *testing.T) {
	schemaJSON := `{
		"type": "object",
		"properties": {
			"enabled": { "type": "boolean" },
			"max_limit": { "type": "integer" }
		},
		"required": ["enabled", "max_limit"]
	}`

	inputJSON := `{
		"enabled": true
	}`

	if isValidInput(schemaJSON, inputJSON) {
		t.Fatal("expected validation to fail for missing required property, but it passed")
	}
}

func TestIsValidInput_InvalidType(t *testing.T) {
	schemaJSON := `{
		"type": "object",
		"properties": {
			"max_limit": { "type": "integer" }
		},
		"required": ["max_limit"]
	}`

	inputJSON := `{
		"max_limit": "not-an-integer"
	}`

	if isValidInput(schemaJSON, inputJSON) {
		t.Fatal("expected validation to fail for invalid type, but it passed")
	}
}

func TestIsValidInput_InvalidJSON(t *testing.T) {
	schemaJSON := `{
		"type": "object",
		"properties": {
			"max_limit": { "type": "integer" }
		},
		"required": ["max_limit"]
	}`

	// Broken JSON
	inputJSON := `{
		"max_limit": 100,,
	}`

	if isValidInput(schemaJSON, inputJSON) {
		t.Fatal("expected validation to fail for invalid JSON, but it passed")
	}
}
