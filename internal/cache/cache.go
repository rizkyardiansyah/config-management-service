package cache

import (
	"sync"

	"sass.com/configsvc/internal/models"
)

var (
	// Mimicking Redis cache
	cache *sync.Map
)

func Init() {
	cache = &sync.Map{}
}

func Put(name string, lastCfg *models.LastConfigurations) {
	cache.Store(name, lastCfg)
}

func Get(name string) (*models.LastConfigurations, bool) {
	val, ok := cache.Load(name)
	if !ok {
		return nil, false
	}
	return val.(*models.LastConfigurations), true
}

func Remove(name string) {
	cache.Delete(name)
}
