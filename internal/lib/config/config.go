package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Env string `env:"ENV" env-default:"local"`
	HTTPServer
	PG
	API
}

type HTTPServer struct {
	Host        string        `env:"SERVER_HOST" env-default:"localhost"`
	Port        int           `env:"SERVER_PORT" env-default:"8080"`
	Timeout     time.Duration `env:"TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-default:"30s"`
}

type PG struct {
	Name           string `env:"PG_NAME" env-required:"true"`
	User           string `env:"PG_USER" env-required:"true"`
	Password       string `env:"PG_PASSWORD" env-required:"true"`
	Host           string `env:"PG_HOST" env-required:"true"`
	Port           string `env:"PG_PORT" env-required:"true"`
	MigrationsPath string `env:"MIGRATIONS_PATH" env-required:"true"`
}

type API struct {
	PersonPageLimit int `env:"PERSON_PAGE_LIMIT" env-default:"5"`
}

func Load() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

type TestConfig struct {
	PG
}

func LoadTest() (*TestConfig, error) {
	var cfg TestConfig

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
