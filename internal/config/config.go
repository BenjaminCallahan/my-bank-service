package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	DBName string `env:"DB_NAME" env-default:"bank.db"`
	Port   string `env:"PORT" env-default:"8080"`
}

func Init() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return &cfg, err
	}
	return &cfg, nil
}
