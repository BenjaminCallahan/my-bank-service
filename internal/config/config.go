package config

import "github.com/ilyakaznacheev/cleanenv"

// Config represent entity of application configuration
type Config struct {
	// Path to the database file
	DBName string `env:"DB_NAME" env-default:"./sqlite_database/bank.db"`
	// Port for web application
	Port string `env:"PORT" env-default:"8080"`
}

// Init initializing a new Config
func Init() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return &cfg, err
	}
	return &cfg, nil
}
