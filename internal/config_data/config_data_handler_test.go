package configdata

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"sass.com/configsvc/internal/models"
)

type mockConfigService struct {
	createErr   error
	updateErr   error
	rollbackErr error
	lastCfg     *models.LastConfigurations
	lastErr     error
	byVerCfg    *models.Configurations
	byVerErr    error
	versions    []models.Configurations
	versionsErr error
}

func (m *mockConfigService) Create(cfg *models.Configurations) error {
	return m.createErr
}
func (m *mockConfigService) Update(cfg *models.Configurations) error {
	return m.updateErr
}
func (m *mockConfigService) RollbackConfig(cfg *models.Configurations) error {
	return m.rollbackErr
}
func (m *mockConfigService) GetLastVersionByName(name string) (*models.LastConfigurations, error) {
	return m.lastCfg, m.lastErr
}
func (m *mockConfigService) GetByNameByVersion(name string, version int) (*models.Configurations, error) {
	return m.byVerCfg, m.byVerErr
}
func (m *mockConfigService) GetConfigVersions(name string) ([]models.Configurations, error) {
	return m.versions, m.versionsErr
}

func setupGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestConfigHandler_CreateConfig_Success(t *testing.T) {
	svc := &mockConfigService{}
	h := NewConfigHandler(svc)
	r := setupGin()
	r.POST("/configs", func(c *gin.Context) {
		c.Set("role", "admin")
		c.Set("user_id", "tester")
		h.CreateConfig(c)
	})

	body := bytes.NewBufferString(`{
		"name":"feature_flag",
		"type":"object",
		"schema":"{\"type\":\"object\",\"properties\":{\"enabled\":{\"type\":\"boolean\"}},\"required\":[\"enabled\"]}",
		"input":"{\"enabled\":true}"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/configs", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Result().StatusCode)
	}
}

func TestConfigHandler_CreateConfig_InvalidBody(t *testing.T) {
	svc := &mockConfigService{}
	h := NewConfigHandler(svc)
	r := setupGin()
	r.POST("/configs", h.CreateConfig)

	body := bytes.NewBufferString(`not-a-json`)
	req := httptest.NewRequest(http.MethodPost, "/configs", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Result().StatusCode)
	}
}

func TestConfigHandler_UpdateConfig_InvalidID(t *testing.T) {
	svc := &mockConfigService{}
	h := NewConfigHandler(svc)
	r := setupGin()
	r.PUT("/configs/:id", h.UpdateConfig)

	body := bytes.NewBufferString(`{}`)
	req := httptest.NewRequest(http.MethodPut, "/configs/not-a-uuid", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Result().StatusCode)
	}
}

func TestConfigHandler_GetConfigByNameByVersion_NotFound(t *testing.T) {
	svc := &mockConfigService{byVerErr: errors.New("not found")}
	h := NewConfigHandler(svc)
	r := setupGin()
	r.GET("/configs/:name/:version", h.GetConfigByNameByVersion)

	req := httptest.NewRequest(http.MethodGet, "/configs/feature_flag/99", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Result().StatusCode)
	}
}

func TestConfigHandler_GetConfigVersions_Error(t *testing.T) {
	svc := &mockConfigService{versionsErr: errors.New("db error")}
	h := NewConfigHandler(svc)
	r := setupGin()
	r.GET("/configs/:name/versions", h.GetConfigVersions)

	req := httptest.NewRequest(http.MethodGet, "/configs/feature_flag/versions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Result().StatusCode)
	}
}

func TestConfigHandler_CreateConfig_AdminSuccess(t *testing.T) {
	svc := &mockConfigService{}
	h := NewConfigHandler(svc)
	r := setupGin()
	r.POST("/configs", func(c *gin.Context) {
		c.Set("role", "admin")
		c.Set("user_id", "tester")
		h.CreateConfig(c)
	})

	body := bytes.NewBufferString(`{
		"name":"feature_flag",
		"type":"object",
		"schema":"{\"type\":\"object\",\"properties\":{\"enabled\":{\"type\":\"boolean\"}},\"required\":[\"enabled\"]}",
		"input":"{\"enabled\":true}"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/configs", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Result().StatusCode)
	}
}

func TestConfigHandler_CreateConfig_UserFailed(t *testing.T) {
	svc := &mockConfigService{}
	h := NewConfigHandler(svc)
	r := setupGin()
	r.POST("/configs", func(c *gin.Context) {
		c.Set("role", "user") // 👈 non-admin
		c.Set("user_id", "tester")
		h.CreateConfig(c)
	})

	body := bytes.NewBufferString(`{
		"name":"feature_flag",
		"type":"object",
		"schema":"{\"type\":\"object\",\"properties\":{\"enabled\":{\"type\":\"boolean\"}},\"required\":[\"enabled\"]}",
		"input":"{\"enabled\":true}"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/configs", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Result().StatusCode)
	}
}

func TestConfigHandler_UpdateConfig_AdminSuccess(t *testing.T) {
	svc := &mockConfigService{
		lastCfg: &models.LastConfigurations{
			Name:    "feature_flag",
			Schema:  `{"type":"object","properties":{"enabled":{"type":"boolean"}},"required":["enabled"]}`,
			Version: 1,
		},
	}
	h := NewConfigHandler(svc)
	r := setupGin()
	r.PUT("/configs/:name", func(c *gin.Context) {
		c.Set("role", "admin")
		c.Set("user_id", "tester")
		h.UpdateConfig(c)
	})

	body := bytes.NewBufferString(`{
		"schema":"{\"type\":\"object\",\"properties\":{\"enabled\":{\"type\":\"boolean\"}},\"required\":[\"enabled\"]}",
		"input":"{\"enabled\":false}"
	}`)

	req := httptest.NewRequest(http.MethodPut, "/configs/feature_flag", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Result().StatusCode)
	}
}

func TestConfigHandler_UpdateConfig_UserUnauthorized(t *testing.T) {
	svc := &mockConfigService{
		lastCfg: &models.LastConfigurations{
			Name:    "feature_flag",
			Schema:  `{"type":"object","properties":{"enabled":{"type":"boolean"}},"required":["enabled"]}`,
			Version: 1,
		},
	}
	h := NewConfigHandler(svc)
	r := setupGin()
	r.PUT("/configs/:name", func(c *gin.Context) {
		c.Set("role", "user") // non-admin
		c.Set("user_id", "tester")
		h.UpdateConfig(c)
	})

	body := bytes.NewBufferString(`{
		"schema":"{\"type\":\"object\",\"properties\":{\"enabled\":{\"type\":\"boolean\"}},\"required\":[\"enabled\"]}",
		"input":"{\"enabled\":true}"
	}`)

	req := httptest.NewRequest(http.MethodPut, "/configs/feature_flag", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Result().StatusCode)
	}
}

func TestConfigHandler_UpdateConfig_SchemaChanged(t *testing.T) {
	svc := &mockConfigService{
		lastCfg: &models.LastConfigurations{
			Name:    "feature_flag",
			Schema:  `{"type":"object","properties":{"enabled":{"type":"boolean"}},"required":["enabled"]}`,
			Version: 1,
		},
	}
	h := NewConfigHandler(svc)
	r := setupGin()
	r.PUT("/configs/:name", func(c *gin.Context) {
		c.Set("role", "admin")
		c.Set("user_id", "tester")
		h.UpdateConfig(c)
	})

	// schema changed (max_limit vs enabled)
	body := bytes.NewBufferString(`{
		"schema":"{\"type\":\"object\",\"properties\":{\"max_limit\":{\"type\":\"integer\"}},\"required\":[\"max_limit\"]}",
		"input":"{\"max_limit\":1000}"
	}`)

	req := httptest.NewRequest(http.MethodPut, "/configs/feature_flag", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Result().StatusCode)
	}
}
