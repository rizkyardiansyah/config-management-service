package auth

import (
	"testing"

	"github.com/google/uuid"
	"sass.com/configsvc/internal/config"
	"sass.com/configsvc/internal/models"
	"sass.com/configsvc/internal/secrets"
)

type mockUserRepo struct {
	user *models.User
}

func (m *mockUserRepo) FindByUsername(username string) (*models.User, error) {
	return m.user, nil
}

func (m *mockUserRepo) VerifyPassword(u *models.User, password string) bool {
	return password == "correct-password"
}

func TestAuthService_Login_Success(t *testing.T) {
	user := &models.User{ID: uuid.New(), Username: "elon", Role: models.RoleUser}

	fakeCfg := &config.Config{AccessTokenTTLInMinutes: 1}
	fakeSecrets := &secrets.Secrets{JWTsecret: []byte("testsecret")}

	svc := NewAuthService(&mockUserRepo{user: user}, fakeCfg, fakeSecrets)

	access, refresh, err := svc.Login("elon", "correct-password")
	if err != nil {
		t.Fatalf("expected success, got err: %v", err)
	}
	if access == "" || refresh == "" {
		t.Fatal("expected non-empty tokens")
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	user := &models.User{ID: uuid.New(), Username: "elon", Role: models.RoleUser}

	fakeCfg := &config.Config{AccessTokenTTLInMinutes: 1}
	fakeSecrets := &secrets.Secrets{JWTsecret: []byte("testsecret")}

	svc := NewAuthService(&mockUserRepo{user: user}, fakeCfg, fakeSecrets)

	_, _, err := svc.Login("elon", "wrong-password")
	if err == nil {
		t.Fatal("expected error for invalid password")
	}
}
