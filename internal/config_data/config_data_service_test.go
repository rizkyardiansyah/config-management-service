package configdata

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"sass.com/configsvc/internal/models"
)

type mockConfigRepo struct {
	createErr    error
	updateErr    error
	getByNameCfg *models.Configurations
	getByNameErr error
	getByVerCfg  *models.Configurations
	getByVerErr  error
	getVersCfgs  []models.Configurations
	getVersErr   error
}

func (m *mockConfigRepo) Create(cfg *models.Configurations) error {
	return m.createErr
}
func (m *mockConfigRepo) Update(cfg *models.Configurations) error {
	return m.updateErr
}
func (m *mockConfigRepo) GetLastVersionByName(name string) (*models.Configurations, error) {
	return m.getByNameCfg, m.getByNameErr
}
func (m *mockConfigRepo) GetByNameByVersion(name string, version int) (*models.Configurations, error) {
	return m.getByVerCfg, m.getByVerErr
}
func (m *mockConfigRepo) GetConfigVersions(name string) ([]models.Configurations, error) {
	return m.getVersCfgs, m.getVersErr
}

func TestConfigService_Create_Success(t *testing.T) {
	mockRepo := &mockConfigRepo{}
	svc := &ConfigServiceImpl{repo: mockRepo}

	cfg := &models.Configurations{ID: uuid.New(), Name: "feature_flag", Version: 1}
	if err := svc.Create(cfg); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestConfigService_Create_Error(t *testing.T) {
	mockRepo := &mockConfigRepo{createErr: errors.New("create failed")}
	svc := &ConfigServiceImpl{repo: mockRepo}

	err := svc.Create(&models.Configurations{})
	if err == nil || err.Error() != "create failed" {
		t.Fatalf("expected 'create failed', got %v", err)
	}
}

func TestConfigService_Update_Error(t *testing.T) {
	mockRepo := &mockConfigRepo{updateErr: errors.New("update failed")}
	svc := &ConfigServiceImpl{repo: mockRepo}

	err := svc.Update(&models.Configurations{})
	if err == nil || err.Error() != "update failed" {
		t.Fatalf("expected 'update failed', got %v", err)
	}
}

func TestConfigService_RollbackConfig(t *testing.T) {
	mockRepo := &mockConfigRepo{}
	svc := &ConfigServiceImpl{repo: mockRepo}

	cfg := &models.Configurations{
		ID:        uuid.New(),
		Name:      "feature_flag",
		Version:   1,
		CreatedAt: time.Now(),
	}

	if err := svc.RollbackConfig(cfg); err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestConfigService_GetByName_Success(t *testing.T) {
	expected := &models.Configurations{Name: "feature_flag", Version: 2}
	mockRepo := &mockConfigRepo{getByNameCfg: expected}
	svc := &ConfigServiceImpl{repo: mockRepo}

	cfg, err := svc.GetLastVersionByName("feature_flag")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if cfg.Version != 2 {
		t.Errorf("expected version 2, got %d", cfg.Version)
	}
}

func TestConfigService_GetByName_Error(t *testing.T) {
	mockRepo := &mockConfigRepo{getByNameErr: errors.New("not found")}
	svc := &ConfigServiceImpl{repo: mockRepo}

	_, err := svc.GetLastVersionByName("missing")
	if err == nil || err.Error() != "not found" {
		t.Fatalf("expected 'not found', got %v", err)
	}
}

func TestConfigService_GetByNameByVersion_Success(t *testing.T) {
	expected := &models.Configurations{Name: "feature_flag", Version: 1}
	mockRepo := &mockConfigRepo{getByVerCfg: expected}
	svc := &ConfigServiceImpl{repo: mockRepo}

	cfg, err := svc.GetByNameByVersion("feature_flag", 1)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if cfg.Version != 1 {
		t.Errorf("expected version 1, got %d", cfg.Version)
	}
}

func TestConfigService_GetByNameByVersion_Error(t *testing.T) {
	mockRepo := &mockConfigRepo{getByVerErr: errors.New("not found")}
	svc := &ConfigServiceImpl{repo: mockRepo}

	_, err := svc.GetByNameByVersion("feature_flag", 99)
	if err == nil || err.Error() != "not found" {
		t.Fatalf("expected 'not found', got %v", err)
	}
}

func TestConfigService_GetConfigVersions_Success(t *testing.T) {
	expected := []models.Configurations{
		{Name: "feature_flag", Version: 1},
		{Name: "feature_flag", Version: 2},
	}
	mockRepo := &mockConfigRepo{getVersCfgs: expected}
	svc := &ConfigServiceImpl{repo: mockRepo}

	cfgs, err := svc.GetConfigVersions("feature_flag")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(cfgs) != 2 {
		t.Errorf("expected 2 configs, got %d", len(cfgs))
	}
}

func TestConfigService_GetConfigVersions_Error(t *testing.T) {
	mockRepo := &mockConfigRepo{getVersErr: errors.New("db error")}
	svc := &ConfigServiceImpl{repo: mockRepo}

	_, err := svc.GetConfigVersions("feature_flag")
	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected 'db error', got %v", err)
	}
}
