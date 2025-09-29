package auth

import (
	"testing"

	"github.com/google/uuid"
	"sass.com/configsvc/internal/config"
	"sass.com/configsvc/internal/models"
	"sass.com/configsvc/internal/secrets"
)

func TestCreateAccessTokenAndParse(t *testing.T) {
	u := models.User{ID: uuid.New(), Role: models.RoleUser}

	fakeCfg := &config.Config{AccessTokenTTLInMinutes: 1}
	fakeSecrets := &secrets.Secrets{JWTsecret: []byte("testsecret")}

	token, err := createAccessToken(u, fakeCfg, fakeSecrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatalf("unexpected empty token")
	}
}

func TestCreateRefreshTokenUnique(t *testing.T) {
	token1, err := createRefreshToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	token2, err := createRefreshToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token1 == token2 {
		t.Fatal("refresh token is not unique")
	}
}
