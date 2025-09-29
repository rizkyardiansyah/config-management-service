package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sass.com/configsvc/internal/auth"
	"sass.com/configsvc/internal/config"
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

	// Setup routes
	r := gin.Default()
	r.POST("api/login", func(c *gin.Context) {
		authHandler.Login(c.Writer, c.Request)
	})

	// Run server using port from config
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("server running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
