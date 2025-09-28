package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// TODO: store this in OS Env
	//secret := []byte("supersecretkey")

	// Public route to get a token
	r.POST("/login", func(c *gin.Context) {
		// var body struct {
		// 	Username string   `json:"username"`
		// 	Roles    []string `json:"roles`
		// }
		c.JSON(http.StatusOK, gin.H{"token": "ddsds"})
	})

	r.Run(":8087")

}
