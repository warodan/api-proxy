package config

import (
	"errors"
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

func Load(tempLogger *slog.Logger) *Config {
	config := &Config{}

	data, err := os.ReadFile("config.yaml")
	if err == nil {
		if yamlErr := yaml.Unmarshal(data, config); yamlErr != nil {
			panic("YAML parsing error: " + yamlErr.Error())
		}
	}

	if err := godotenv.Load(); err != nil {
		tempLogger.Error("No .env file found, using system env vars")
	}

	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
	}
	if level := os.Getenv("LOGGER_LEVEL"); level != "" {
		config.LoggerLevel = level
	}

	if config.Port == "" || config.LoggerLevel == "" {
		panic(errors.New("config is incomplete: port or loggerLevel missing"))
	}

	return config
}
