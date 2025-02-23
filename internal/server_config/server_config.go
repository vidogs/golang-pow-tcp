package server_config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type (
	ServerConfig struct {
		Bind     Bind     `yaml:"bind"`
		Logger   Logger   `yaml:"logger"`
		Settings Settings `yaml:"settings"`
		Dao      Dao      `yaml:"dao"`
	}

	Bind struct {
		Host string `yaml:"host" env:"POW_BIND_HOST"`
		Port int    `yaml:"port" env:"POW_BIND_PORT"`
	}

	Logger struct {
		Level  string `yaml:"level"         env:"POW_LOG_LEVEL"`
		Format string `yaml:"format"        env:"POW_LOG_FORMAT"`
	}

	Settings struct {
		Challenge    Challenge     `yaml:"challenge"`
		SolveTimeout time.Duration `yaml:"solve_timeout" env:"POW_SOLVE_TIMEOUT"`
	}

	Challenge struct {
		Length     int `yaml:"length"     env:"POW_SETTINGS_CHALLENGE_LENGTH"`
		Difficulty int `yaml:"difficulty" env:"POW_SETTINGS_CHALLENGE_DIFFICULTY"`
	}

	Dao struct {
		WordsOfWisdom []string `yaml:"words_of_wisdom"`
	}
)

func NewServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}

	configFileName := os.Getenv("POW_CONFIG_FILE")

	if len(configFileName) == 0 {
		configFileName = "./config/server/config.yaml"
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
