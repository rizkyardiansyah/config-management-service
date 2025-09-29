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

func (r *ConfigRepo) Update(cfg *models.Configurations) error {
	return r.db.Save(cfg).Error
}

// Get latest version of a config by name
func (r *ConfigRepo) GetLastVersionByName(name string) (*models.Configurations, error) {
	var cfg models.Configurations
	if err := r.db.Where("name = ?", name).
		Order("version DESC").
		First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *ConfigRepo) GetByNameByVersion(name string, version int) (*models.Configurations, error) {
	var cfg models.Configurations
	if err := r.db.Where("name = ? AND version = ?", name, version).
		First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *ConfigRepo) GetConfigVersions(name string) ([]models.Configurations, error) {
	var configs []models.Configurations
	if err := r.db.Where("name = ?", name).
		Order("version ASC"). // change to DESC if you prefer latest first
		Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}
