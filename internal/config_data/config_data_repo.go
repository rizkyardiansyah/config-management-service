package configdata

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sass.com/configsvc/internal/models"
)

type ConfigRepo struct {
	db *gorm.DB
}

func NewConfigRepo(db *gorm.DB) *ConfigRepo {
	return &ConfigRepo{db: db}
}

func (r *ConfigRepo) Create(cfg *models.Configurations) error {
	cfg.ID = uuid.New()
	return r.db.Create(cfg).Error
}
