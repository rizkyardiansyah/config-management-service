package configdata

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sass.com/configsvc/internal/models"
)

func setupConfigTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite in-memory: %v", err)
	}
	if err := db.AutoMigrate(&models.Configurations{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}
	return db
}

func TestConfigDataRepo_CreateMultipleProperties_Success(t *testing.T) {
	db := setupConfigTestDB(t)
	repo := NewConfigRepo(db)

	schemaJSON := `{
		"type": "object",
		"properties": {
			"max_limit": { "type": "integer" },
			"enabled":   { "type": "boolean" }
		},
		"required": ["max_limit", "enabled"]
	}`

	inputJSON := `{
		"max_limit": 100000,
		"enabled": true
	}`

	if !isValidInput(schemaJSON, inputJSON) {
		t.Fatal("Schema and input are conflicting")
	}

	cfg := &models.Configurations{
		ID:        uuid.New(),
		ClientID:  "bca-pusat",
		Name:      "BCA_VA_DAILY_TRESHOLD",
		Type:      models.TypeObject,
		Schema:    schemaJSON,
		Input:     inputJSON,
		Version:   1,
		CreatedAt: time.Now(),
		CreatedBy: "herman@bca.id",
		UpdatedAt: time.Now(),
		IsActive:  1,
	}

	if err := repo.Create(cfg); err != nil {
		t.Fatalf("failed to create config: %v", err)
	}
}

func TestConfigDataRepo_CreateSingleProperty_Success(t *testing.T) {
	db := setupConfigTestDB(t)
	repo := NewConfigRepo(db)

	schemaJSON := `{
		"type": "object",
		"properties": {
		  "max_transfer": { "type": "integer" }
		},
		"required": ["max_transfer"]
	  }`

	inputJSON := `{
		"max_transfer": 2333000
	  }`

	if !isValidInput(schemaJSON, inputJSON) {
		t.Fatal("Schema and input not match")
	}

	cfg := &models.Configurations{
		ID:        uuid.New(),
		ClientID:  "bca-pusat",
		Name:      "BCA_VA_DAILY_TRESHOLD",
		Type:      models.TypeObject,
		Schema:    schemaJSON,
		Input:     inputJSON,
		Version:   1,
		CreatedAt: time.Now(),
		CreatedBy: "herman@bca.id",
		UpdatedAt: time.Now(),
		IsActive:  1,
	}

	if err := repo.Create(cfg); err != nil {
		t.Fatalf("failed to create config: %v", err)
	}
}

func TestConfigDataRepo_CreateMultipleProperties_Failed(t *testing.T) {
	schemaJSON := `{
		"type": "object",
		"properties": {
			"max_limit": { "type": "integer" },
			"enabled":   { "type": "boolean" }
		},
		"required": ["max_limit", "enabled"]
	}`

	inputJSON := `{
		"max_limit": 100000,
	}`

	if isValidInput(schemaJSON, inputJSON) {
		t.Fatal("Missing input proporty, should failed")
	}
}

func TestConfigDataRepo_CreateSingleProperty_Failed(t *testing.T) {
	schemaJSON := `{
		"type": "object",
		"properties": {
		  "max_transfer": { "type": "integer" }
		},
		"required": ["max_transfer"]
	  }`

	inputJSON := `{
		"max_transfer": "2333000"
	  }`

	if isValidInput(schemaJSON, inputJSON) {
		t.Fatal("Input property is different with schema, should failed")
	}
}
