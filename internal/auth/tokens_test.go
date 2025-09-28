package auth

import (
	"testing"

	"github.com/google/uuid"
	"sass.com/configsvc/internal/models"
)

func TestCreateAccessTokenAndParse(t *testing.T) {
	u := models.User{ID: uuid.New(), Role: models.RoleUser}
	token, err := createAccessToken(u)
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
