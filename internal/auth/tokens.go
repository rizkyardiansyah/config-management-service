package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"sass.com/configsvc/internal/models"
)

// TODO: load this from config and Env
var (
	jwtSecret       = []byte("replace-with-secure-secret")
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

func createAccessToken(u models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  u.ID,
		"role": string(u.Role),
		"exp":  time.Now().Add(accessTokenTTL).Unix(),
		"iat":  time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret)
}
