package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log/slog"
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
	var yamlErr error
	var envLoaded bool

	if data, err := os.ReadFile("config.yaml"); err == nil {
		if err := yaml.Unmarshal(data, config); err != nil {
			yamlErr = fmt.Errorf("failed to parse config.yaml: %w", err)
		}
	} else if !os.IsNotExist(err) {
		yamlErr = fmt.Errorf("failed to read config.yaml: %w", err)
	}

	_ = godotenv.Load()

	if val := os.Getenv("PORT"); val != "" {
		config.Port = val
		envLoaded = true
	}
	if val := os.Getenv("LOGGER_LEVEL"); val != "" {
		config.LoggerLevel = val
		envLoaded = true
	}

	if yamlErr != nil && envLoaded {
		slog.Warn("yaml config invalid, using fallback from env", "err", yamlErr)
	}

	if yamlErr == nil && !envLoaded {
		slog.Warn("env variables not set, using only yaml config")
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return config, nil
}
