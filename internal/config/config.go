package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port                     int
	AccessTokenTTLInDays     int
	RefreshTokenTTLInMinutes int
}

func LoadConfig() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config/config.json"
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	return config, err
}
