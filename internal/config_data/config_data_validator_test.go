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

func TestEqualSchemas_SameFormatting(t *testing.T) {
	a := `{"type":"object","properties":{"enabled":{"type":"boolean"}},"required":["enabled"]}`
	b := `{"type":"object","properties":{"enabled":{"type":"boolean"}},"required":["enabled"]}`

	if !equalSchemas(a, b) {
		t.Fatal("expected schemas to be equal")
	}
}

func TestEqualSchemas_DifferentFormatting(t *testing.T) {
	a := `{
		"type": "object",
		"properties": {
			"enabled": { "type": "boolean" }
		},
		"required": ["enabled"]
	}`
	b := `{"required":["enabled"],"properties":{"enabled":{"type":"boolean"}},"type":"object"}`

	if !equalSchemas(a, b) {
		t.Fatal("expected schemas to be equal despite formatting differences")
	}
}

func TestEqualSchemas_DifferentSchema(t *testing.T) {
	a := `{"type":"object","properties":{"enabled":{"type":"boolean"}},"required":["enabled"]}`
	b := `{"type":"object","properties":{"max_limit":{"type":"integer"}},"required":["max_limit"]}`

	if equalSchemas(a, b) {
		t.Fatal("expected schemas to be different")
	}
}

func TestEqualSchemas_InvalidJSON(t *testing.T) {
	a := `{"type":"object"`
	b := `{"type":"object"}`

	if equalSchemas(a, b) {
		t.Fatal("expected invalid JSON not to be equal")
	}
}
