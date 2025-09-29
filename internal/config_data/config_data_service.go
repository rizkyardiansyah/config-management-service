package configdata

import (
	"sass.com/configsvc/internal/models"
)

type ConfigService interface {
	Create(cfg *models.Configurations) error
	Update(cfg *models.Configurations) error
	RollbackConfig(cfg *models.Configurations) error
	GetByName(name string) (*models.Configurations, error)
	GetByNameByVersion(name string, version int) (*models.Configurations, error)
	GetConfigVersions(name string) ([]models.Configurations, error)
}

type ConfigServiceImpl struct {
	repo ConfigRepository
}

type ConfigRepository interface {
	Create(cfg *models.Configurations) error
	Update(cfg *models.Configurations) error
	GetByName(name string) (*models.Configurations, error)
	GetByNameByVersion(name string, version int) (*models.Configurations, error)
	GetConfigVersions(name string) ([]models.Configurations, error)
}

func NewConfigService(repo *ConfigRepo) ConfigService {
	return &ConfigServiceImpl{repo: repo}
}

func (s *ConfigServiceImpl) Create(cfg *models.Configurations) error {
	return s.repo.Create(cfg)
}

func (s *ConfigServiceImpl) Update(cfg *models.Configurations) error {
	return s.repo.Update(cfg)
}

func (s *ConfigServiceImpl) RollbackConfig(cfg *models.Configurations) error {
	return s.repo.Create(cfg)
}

func (s *ConfigServiceImpl) GetByName(name string) (*models.Configurations, error) {
	return s.repo.GetByName(name)
}

func (s *ConfigServiceImpl) GetByNameByVersion(name string, version int) (*models.Configurations, error) {
	return s.repo.GetByNameByVersion(name, version)
}

func (s *ConfigServiceImpl) GetConfigVersions(name string) ([]models.Configurations, error) {
	return s.repo.GetConfigVersions(name)
}
