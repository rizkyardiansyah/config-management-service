package auth

import (
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sass.com/configsvc/internal/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}
	return db
}

func TestUserRepo_FindByUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepo(db)

	pw, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	user := models.User{ID: uuid.New(), Username: "alice", PasswordHash: string(pw), Role: models.RoleUser}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	found, err := repo.FindByUsername("alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil || found.Username != "alice" {
		t.Fatalf("expected alice, got %+v", found)
	}

	notFound, err := repo.FindByUsername("bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if notFound != nil {
		t.Fatalf("expected nil for bob, got %+v", notFound)
	}
}

func TestUserRepo_VerifyPassword(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepo(db)

	pw, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	user := models.User{ID: uuid.New(), Username: "alice", PasswordHash: string(pw), Role: models.RoleUser}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	ok := repo.VerifyPassword(&user, "secret")
	if !ok {
		t.Fatal("expected password to verify")
	}

	notOk := repo.VerifyPassword(&user, "wrong")
	if notOk {
		t.Fatal("expected password to fail")
	}
}
