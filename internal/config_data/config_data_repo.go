package configdata

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sass.com/configsvc/internal/models"
)

type ConfigRepo interface {
	Create(cfg *models.Configurations, lastCfg *models.LastConfigurations) error
	Update(cfg *models.Configurations) error
	GetLastConfig(name string) (*models.LastConfigurations, error)
	GetByNameByVersion(name string, version int) (*models.Configurations, error)
	GetConfigVersions(name string) ([]models.Configurations, error)
}

func NewConfigRepo(db *gorm.DB) ConfigRepo {
	return &ConfigRepoImpl{db: db}
}

type ConfigRepoImpl struct {
	db *gorm.DB
}

func (r *ConfigRepoImpl) Create(cfg *models.Configurations, last *models.LastConfigurations) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(cfg).Error; err != nil {
			return err
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{"client_id", "type", "schema", "input", "version", "updated_at", "is_active"}),
		}).Create(last).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *ConfigRepoImpl) Update(cfg *models.Configurations) error {
	return r.db.Save(cfg).Error
}

// Get latest version of a config by name
func (r *ConfigRepoImpl) GetLastConfig(name string) (*models.LastConfigurations, error) {
	var lastCfg models.LastConfigurations
	if err := r.db.Where("name = ?", name).
		First(&lastCfg).Error; err != nil {
		return nil, err
	}
	return &lastCfg, nil
}

func (r *ConfigRepoImpl) GetByNameByVersion(name string, version int) (*models.Configurations, error) {
	var cfg models.Configurations
	if err := r.db.Where("name = ? AND version = ?", name, version).
		First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *ConfigRepoImpl) GetConfigVersions(name string) ([]models.Configurations, error) {
	var configs []models.Configurations
	if err := r.db.Where("name = ?", name).
		Order("version ASC"). // change to DESC if you prefer latest first
		Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}
