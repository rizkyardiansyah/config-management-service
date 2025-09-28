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
		t.Fatalf("expected non-empty token")
	}
}
