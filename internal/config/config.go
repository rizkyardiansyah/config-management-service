package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port                     int
	AccessTokenTTLInMinutes  int
	RefreshTokenTTLInMinutes int
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("../../config/config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	return config, err
}
