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
	// migrate both history and last snapshot tables so repo.Create(..., last) works
	if err := db.AutoMigrate(&models.Configurations{}, &models.LastConfigurations{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}
	return db
}

func makeLastFromCfg(cfg *models.Configurations) *models.LastConfigurations {
	return &models.LastConfigurations{
		ID:        uuid.New(),
		ClientID:  cfg.ClientID,
		Name:      cfg.Name,
		Type:      cfg.Type,
		Schema:    cfg.Schema,
		Input:     cfg.Input,
		Version:   cfg.Version,
		CreatedAt: cfg.CreatedAt,
		CreatedBy: cfg.CreatedBy,
		UpdatedAt: cfg.UpdatedAt,
		IsActive:  cfg.IsActive,
	}
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
		t.Fatal("validation failed: schema and input are conflicting")
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

	// create last snapshot from cfg and pass to Create
	if err := repo.Create(cfg, makeLastFromCfg(cfg)); err != nil {
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
		t.Fatal("validation failed: schema and input are conflicting")
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

	if err := repo.Create(cfg, makeLastFromCfg(cfg)); err != nil {
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
		"max_limit": 100000
	}`

	if isValidInput(schemaJSON, inputJSON) {
		t.Fatal("validation should fail: missing required property 'enabled'")
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
		t.Fatal("validation should fail: 'max_transfer' type mismatch")
	}
}

func TestConfigDataRepo_Update_Success(t *testing.T) {
	db := setupConfigTestDB(t)
	repo := NewConfigRepo(db)

	schemaJSON := `{
		"type": "object",
		"properties": {
			"enabled": { "type": "boolean" }
		},
		"required": ["enabled"]
	}`

	inputJSON := `{
		"enabled": true
	}`

	if !isValidInput(schemaJSON, inputJSON) {
		t.Fatal("validation failed: schema and input are conflicting")
	}

	cfg := &models.Configurations{
		ID:        uuid.New(),
		ClientID:  "bca-pusat",
		Name:      "feature_flag",
		Type:      models.TypeObject,
		Schema:    schemaJSON,
		Input:     inputJSON,
		Version:   1,
		CreatedAt: time.Now(),
		CreatedBy: "tester",
		UpdatedAt: time.Now(),
		IsActive:  1,
	}
	if err := repo.Create(cfg, makeLastFromCfg(cfg)); err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	newInput := `{
		"enabled": false
	}`

	if !isValidInput(schemaJSON, newInput) {
		t.Fatal("validation failed: schema and new input are conflicting")
	}

	cfg.Input = newInput
	if err := repo.Update(cfg); err != nil {
		t.Fatalf("failed to update config: %v", err)
	}

	var updated models.Configurations
	if err := db.First(&updated, "id = ?", cfg.ID).Error; err != nil {
		t.Fatalf("failed to fetch updated config: %v", err)
	}
	if updated.Input != newInput {
		t.Errorf("expected updated input %s, got %s", newInput, updated.Input)
	}
}

func TestConfigDataRepo_GetByName(t *testing.T) {
	db := setupConfigTestDB(t)
	repo := NewConfigRepo(db)

	schemaJSON := `{
		"type": "object",
		"properties": {
			"v": { "type": "integer" }
		},
		"required": ["v"]
	}`

	inputV1 := `{
		"v": 1
	}`

	inputV2 := `{
		"v": 2
	}`

	if !isValidInput(schemaJSON, inputV1) || !isValidInput(schemaJSON, inputV2) {
		t.Fatal("validation failed: schema and input are conflicting")
	}

	cfg1 := &models.Configurations{ID: uuid.New(), Name: "feature_flag", Version: 1, Schema: schemaJSON, Input: inputV1}
	cfg2 := &models.Configurations{ID: uuid.New(), Name: "feature_flag", Version: 2, Schema: schemaJSON, Input: inputV2}
	_ = repo.Create(cfg1, makeLastFromCfg(cfg1))
	_ = repo.Create(cfg2, makeLastFromCfg(cfg2))

	latest, err := repo.GetLastConfig("feature_flag")
	if err != nil {
		t.Fatalf("failed to get config by name: %v", err)
	}
	if latest.Version != 2 {
		t.Errorf("expected latest version 2, got %d", latest.Version)
	}
}

func TestConfigDataRepo_GetByName_NotFound(t *testing.T) {
	db := setupConfigTestDB(t)
	repo := NewConfigRepo(db)

	_, err := repo.GetLastConfig("does_not_exist")
	if err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

func TestConfigDataRepo_GetByNameByVersion(t *testing.T) {
	db := setupConfigTestDB(t)
	repo := NewConfigRepo(db)

	schemaJSON := `{
		"type": "object",
		"properties": {
			"v": { "type": "integer" }
		},
		"required": ["v"]
	}`

	inputV1 := `{
		"v": 1
	}`

	inputV2 := `{
		"v": 2
	}`

	if !isValidInput(schemaJSON, inputV1) || !isValidInput(schemaJSON, inputV2) {
		t.Fatal("validation failed: schema and input are conflicting")
	}

	cfg1 := &models.Configurations{ID: uuid.New(), Name: "feature_flag", Version: 1, Schema: schemaJSON, Input: inputV1}
	cfg2 := &models.Configurations{ID: uuid.New(), Name: "feature_flag", Version: 2, Schema: schemaJSON, Input: inputV2}
	_ = repo.Create(cfg1, makeLastFromCfg(cfg1))
	_ = repo.Create(cfg2, makeLastFromCfg(cfg2))

	v1, err := repo.GetByNameByVersion("feature_flag", 1)
	if err != nil {
		t.Fatalf("failed to get config by name and version: %v", err)
	}
	if v1.Version != 1 {
		t.Errorf("expected version 1, got %d", v1.Version)
	}
}

func TestConfigDataRepo_GetByNameByVersion_NotFound(t *testing.T) {
	db := setupConfigTestDB(t)
	repo := NewConfigRepo(db)

	_, err := repo.GetByNameByVersion("missing", 99)
	if err == nil {
		t.Fatal("expected error for missing config version, got nil")
	}
}

func TestConfigDataRepo_GetConfigVersions(t *testing.T) {
	db := setupConfigTestDB(t)
	repo := NewConfigRepo(db)

	schemaJSON := `{
		"type": "object",
		"properties": {
			"v": { "type": "integer" }
		},
		"required": ["v"]
	}`

	inputV1 := `{
		"v": 1
	}`

	inputV2 := `{
		"v": 2
	}`

	inputV3 := `{
		"v": 3
	}`

	if !isValidInput(schemaJSON, inputV1) || !isValidInput(schemaJSON, inputV2) || !isValidInput(schemaJSON, inputV3) {
		t.Fatal("validation failed: schema and input are conflicting")
	}

	cfg1 := &models.Configurations{ID: uuid.New(), Name: "feature_flag", Version: 1, Schema: schemaJSON, Input: inputV1}
	cfg2 := &models.Configurations{ID: uuid.New(), Name: "feature_flag", Version: 2, Schema: schemaJSON, Input: inputV2}
	cfg3 := &models.Configurations{ID: uuid.New(), Name: "feature_flag", Version: 3, Schema: schemaJSON, Input: inputV3}
	_ = repo.Create(cfg1, makeLastFromCfg(cfg1))
	_ = repo.Create(cfg2, makeLastFromCfg(cfg2))
	_ = repo.Create(cfg3, makeLastFromCfg(cfg3))

	list, err := repo.GetConfigVersions("feature_flag")
	if err != nil {
		t.Fatalf("failed to get config versions: %v", err)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 configs, got %d", len(list))
	}
	if list[0].Version != 1 || list[2].Version != 3 {
		t.Errorf("expected versions [1,2,3], got %+v", list)
	}
}

func TestConfigDataRepo_GetConfigVersions_NotFound(t *testing.T) {
	db := setupConfigTestDB(t)
	repo := NewConfigRepo(db)

	list, err := repo.GetConfigVersions("does_not_exist")
	if err != nil {
		t.Fatalf("expected empty list, got error: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected 0 results, got %d", len(list))
	}
}
