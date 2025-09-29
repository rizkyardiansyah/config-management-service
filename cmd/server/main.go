package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sass.com/configsvc/internal/auth"
	"sass.com/configsvc/internal/config"
	configdata "sass.com/configsvc/internal/config_data"
	"sass.com/configsvc/internal/secrets"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	// Load secrets
	secs := secrets.LoadSecrets()

	// Setup DB
	db, err := gorm.Open(sqlite.Open("./data/config.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Wire repo, service, handler
	userRepo := auth.NewUserRepo(db)
	authService := auth.NewAuthService(userRepo, cfg, secs)
	authHandler := auth.NewAuthHandler(authService)
	configRepo := configdata.NewConfigRepo(db)
	configService := configdata.NewConfigService(configRepo)
	configHandler := configdata.NewConfigHandler(configService)

	// Setup routes
	r := gin.Default()
	r.SetTrustedProxies(nil) // disables trusting any proxy

	// Public routes
	r.POST("/api/v1/login", func(c *gin.Context) {
		authHandler.Login(c.Writer, c.Request)
	})

	// JWT-protected routes
	// Protected group
	api := r.Group("/api/v1")
	api.Use(auth.AuthMiddleware(secs))
	{
		api.POST("/configs", configHandler.CreateConfig)
		api.PUT("/configs/:id", configHandler.UpdateConfig)
		api.POST("/configs/:id/rollback", configHandler.RollbackConfig)
		api.GET("/configs/:name/latest", configHandler.GetLastVersionByName)
		api.GET("/configs/:name/versions/:version", configHandler.GetConfigByNameByVersion)
		api.GET("/configs/:name/versions", configHandler.GetConfigVersions)
	}

	// Run server using port from config
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("server running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
