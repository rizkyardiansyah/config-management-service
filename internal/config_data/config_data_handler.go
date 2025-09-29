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
	var cfg models.Configurations

	fmt.Println("in 1")
	if err := c.ShouldBindJSON(&cfg); err != nil {
		fmt.Println("Bind error:", err)
		fmt.Println("Raw body:", cfg)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	fmt.Println("in 2")

	if !isValidInput(cfg.Schema, cfg.Input) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	fmt.Println("in 3")

	// Assign new UUID if not provided
	if cfg.ID == uuid.Nil {
		cfg.ID = uuid.New()
	}

	cfg.Version = 1

	fmt.Println("in 4")

	if err := h.service.Create(&cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create config"})
		return
	}

	fmt.Println("in 5")
	c.JSON(http.StatusCreated, cfg)
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
