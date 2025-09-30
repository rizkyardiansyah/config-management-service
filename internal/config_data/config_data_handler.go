package configdata

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sass.com/configsvc/internal/models"
)

type ConfigHandler struct {
	service ConfigService
}

func NewConfigHandler(service ConfigService) *ConfigHandler {
	return &ConfigHandler{service: service}
}

func (h *ConfigHandler) CreateConfig(c *gin.Context) {
	var newCfg models.Configurations

	if err := c.ShouldBindJSON(&newCfg); err != nil {
		fmt.Println("Bind error:", err)
		fmt.Println("Raw body:", newCfg)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Reject invalid input and schema pair
	if !isValidInput(newCfg.Schema, newCfg.Input) {
		fmt.Println("Invalid schema input pair")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Only Admin is allowed to create config
	roleVal, exists := c.Get("role")
	if !exists {
		// key not set
		c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found"})
		return
	}
	role, ok := roleVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid role type"})
		return
	}
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you are not authorized"})
		return
	}

	// Enforce user id validation
	userIdVal, exists := c.Get("user_id")
	if !exists {
		fmt.Println("User is not authorized)")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	userId, ok := userIdVal.(string)
	if !ok {
		fmt.Println("User is not authorized)")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	existingCfg, _ := h.service.GetLastVersionByName(newCfg.Name)

	// If config already exist, reject
	if existingCfg != nil {
		fmt.Println("config already exists, please do update instead")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "config already exists"})
		return
	}

	// If config not exist, create new config
	if newCfg.ID == uuid.Nil {
		newCfg.ID = uuid.New()
	}
	newCfg.CreatedBy = userId
	newCfg.Version = 1
	newCfg.IsActive = 1
	if err := h.service.Create(&newCfg); err != nil {
		fmt.Println("service failed to create config")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, newCfg)
}

func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var cfg models.Configurations
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	cfg.ID = id

	if err := h.service.Update(&cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update config"})
		return
	}
	c.JSON(http.StatusOK, cfg)
}

func (h *ConfigHandler) RollbackConfig(c *gin.Context) {
	idParam := c.Param("id")
	_, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var cfg models.Configurations
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.RollbackConfig(&cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rollback config"})
		return
	}
	c.JSON(http.StatusCreated, cfg)
}

func (h *ConfigHandler) GetLastVersionByName(c *gin.Context) {
	name := c.Param("name")
	cfg, err := h.service.GetLastVersionByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "config not found"})
		return
	}
	c.JSON(http.StatusOK, cfg)
}

func (h *ConfigHandler) GetConfigByNameByVersion(c *gin.Context) {
	name := c.Param("name")
	versionStr := c.Param("version")
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid version"})
		return
	}

	cfg, err := h.service.GetByNameByVersion(name, version)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "config version not found"})
		return
	}
	c.JSON(http.StatusOK, cfg)
}

func (h *ConfigHandler) GetConfigVersions(c *gin.Context) {
	name := c.Param("name")
	cfgs, err := h.service.GetConfigVersions(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get config versions"})
		return
	}
	c.JSON(http.StatusOK, cfgs)
}
