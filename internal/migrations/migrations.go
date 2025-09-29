package migrations

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- Models ---
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID           uuid.UUID `gorm:"primarykey"`
	Username     string    `gorm:"uniqueIndex;size:100"`
	PasswordHash string
	Role         Role `gorm:"size:20;default:'user'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// --- Helpers ---
func hashPassword(pw string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

// --- Migration Runner ---
func RunMigrations(dbPath string, withSeed bool) error {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect DB: %w", err)
	}

	// AutoMigrate creates/updates schema
	if err := db.AutoMigrate(&User{}); err != nil {
		return fmt.Errorf("failed to migrate schema: %w", err)
	}
	fmt.Println("schema migrated")

	if withSeed {
		seed(db)
	}

	return nil
}

func seed(db *gorm.DB) {
	admin := User{
		ID:           uuid.New(),
		Username:     "admin",
		PasswordHash: hashPassword("admin123"),
		Role:         RoleAdmin,
	}
	user := User{
		ID:           uuid.New(),
		Username:     "user1",
		PasswordHash: hashPassword("user123"),
		Role:         RoleUser,
	}

	db.FirstOrCreate(&admin, User{Username: admin.Username})
	db.FirstOrCreate(&user, User{Username: user.Username})

	fmt.Println("seeded admin and user1")
}
