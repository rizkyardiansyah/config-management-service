package migrations

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sass.com/configsvc/internal/models"
)

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
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("failed to migrate User schema: %w", err)
	}
	if err := db.AutoMigrate(&models.Configurations{}); err != nil {
		return fmt.Errorf("failed to migrate Configurations schema: %w", err)
	}
	fmt.Println("all schemas migrated")

	if withSeed {
		seed(db)
	}

	return nil
}

func seed(db *gorm.DB) {
	admin := models.User{
		ID:           uuid.New(),
		Username:     "admin",
		PasswordHash: hashPassword("admin123"),
		Role:         models.RoleAdmin,
	}
	user := models.User{
		ID:           uuid.New(),
		Username:     "user1",
		PasswordHash: hashPassword("user123"),
		Role:         models.RoleUser,
	}

	db.FirstOrCreate(&admin, models.User{Username: admin.Username})
	db.FirstOrCreate(&user, models.User{Username: user.Username})

	fmt.Println("seeded admin and user1")
}
