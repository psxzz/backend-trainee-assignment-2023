package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	LogsPath    string `env:"AVITO_LOG_PATH" env-default:"./logs/"`
	DatabaseDSN string `env:"AVITO_DATABASE_DSN" env-default:"postgres://postgres:postgres@localhost:5432/experimental_segments?sslmode=disable"`
}

func New() *Config {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
