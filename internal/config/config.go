package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Port        string `yaml:"port" validate:"required,numeric"`
	LoggerLevel string `yaml:"loggerLevel" validate:"required,oneof=DEBUG INFO WARN ERROR"`
}

func (cfg *Config) Validate() error {
	return validator.New().Struct(cfg)
}

func Load() (*Config, error) {
	config := &Config{}

	data, err := os.ReadFile("config.yaml")
	if err == nil {
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse yaml: %w", err)
		}
	}

	_ = godotenv.Load()

	if val := os.Getenv("PORT"); val != "" {
		config.Port = val
	}
	if val := os.Getenv("LOGGER_LEVEL"); val != "" {
		config.LoggerLevel = val
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return config, nil
}
