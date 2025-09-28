package secrets

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Secrets struct {
	JWTsecret []byte
}

func LoadSecrets() *Secrets {
	// Load .env file for local development
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file, relying on system environment variables")
	}

	return &Secrets{
		JWTsecret: []byte(os.Getenv("JWT_SECRET")),
	}
}
