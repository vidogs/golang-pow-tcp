package client_config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type (
	ClientConfig struct {
		Server Server `yaml:"server"`
		Logger Logger `yaml:"logger"`
	}

	Server struct {
		Host string `yaml:"host" env:"POW_SERVER_HOST"`
		Port int    `yaml:"port" env:"POW_SERVER_PORT"`
	}

	Logger struct {
		Level  string `yaml:"level"         env:"POW_LOG_LEVEL"`
		Format string `yaml:"format"        env:"POW_LOG_FORMAT"`
	}
)

func NewClientConfig() (*ClientConfig, error) {
	cfg := &ClientConfig{}

	configFileName := os.Getenv("POW_CONFIG_FILE")

	if len(configFileName) == 0 {
		configFileName = "./config/client/config.yaml"
	}

	err := cleanenv.ReadConfig(configFileName, cfg)

	if err != nil {
		return nil, fmt.Errorf("error reading config file from %s: %w", configFileName, err)
	}

	if err = cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("error reading env: %w", err)
	}

	return cfg, nil
}
