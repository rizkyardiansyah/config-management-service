package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"sass.com/configsvc/internal/config"
	"sass.com/configsvc/internal/models"
	"sass.com/configsvc/internal/secrets"
)

func createAccessToken(u models.User, cfg *config.Config, secrets *secrets.Secrets) (string, error) {
	claims := jwt.MapClaims{
		"sub":  u.ID,
		"role": string(u.Role),
		"exp":  time.Now().Add(time.Duration(cfg.AccessTokenTTLInDays) * 24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secrets.JWTsecret)
}

func createRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
