package configdata

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"sass.com/configsvc/internal/cache"
	"sass.com/configsvc/internal/models"
)

type ConfigService interface {
	Create(cfg *models.Configurations) error
	Update(cfg *models.Configurations) error
	RollbackConfig(cfg *models.Configurations) error
	GetLastVersionByName(name string) (*models.LastConfigurations, error)
	GetByNameByVersion(name string, version int) (*models.Configurations, error)
	GetConfigVersions(name string) ([]models.Configurations, error)
}

func NewConfigService(repo ConfigRepo) ConfigService {
	return &ConfigServiceImpl{repo: repo}
}

type ConfigServiceImpl struct {
	repo ConfigRepo
}

func (s *ConfigServiceImpl) Create(cfg *models.Configurations) error {
	lastCfg, err := s.GetLastVersionByName(cfg.Name)
	if err != nil {
		return err
	}

	nextVersion := 1
	if lastCfg != nil {
		nextVersion = lastCfg.Version + 1
	}

	cfg.ID = uuid.New()
	cfg.Version = nextVersion

	newLastCfg := &models.LastConfigurations{
		ID:        uuid.New(),
		ClientID:  cfg.ClientID,
		Name:      cfg.Name,
		Type:      cfg.Type,
		Schema:    cfg.Schema,
		Input:     cfg.Input,
		Version:   cfg.Version,
		CreatedBy: cfg.CreatedBy,
		IsActive:  1,
	}

	if err := s.repo.Create(cfg, newLastCfg); err != nil {
		return err
	}

	// Push new data to cache
	cache.Put(cfg.Name, newLastCfg)

	return nil
}

func (s *ConfigServiceImpl) Update(cfg *models.Configurations) error {

	return s.repo.Update(cfg)
}

func (s *ConfigServiceImpl) RollbackConfig(cfg *models.Configurations) error {
	return s.repo.Create(cfg, nil)
}

func (s *ConfigServiceImpl) GetLastVersionByName(name string) (*models.LastConfigurations, error) {
	// Get from cache first
	if cacheData, ok := cache.Get(name); ok && cacheData != nil {
		return cacheData, nil
	}

	dbData, err := s.repo.GetLastConfig(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// First creation just return nil, nil
			return nil, nil
		}
		return nil, err
	}

	if dbData != nil {
		// Push db data to cache
		cache.Put(name, dbData)
	}

	return dbData, nil
}

func (s *ConfigServiceImpl) GetByNameByVersion(name string, version int) (*models.Configurations, error) {
	cfg, err := s.repo.GetByNameByVersion(name, version)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // no such version: return nil, nil
		}
		return nil, err
	}
	return cfg, nil
}

func (s *ConfigServiceImpl) GetConfigVersions(name string) ([]models.Configurations, error) {
	return s.repo.GetConfigVersions(name)
}
