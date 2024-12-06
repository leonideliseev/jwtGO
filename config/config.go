package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const defaultConfigFile = "config.yaml"

type Config struct {
	HTTP       HTTP       `yaml:"http"`
	Postgresql Postgresql `yaml:"postgresql"`
	JWT        JWT        `yaml:"jwt"`
}

type HTTP struct {
	Port string `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	Host string `yaml:"host" env:"HTTP_HOST" env-default:"0.0.0.0"`
}

type Postgresql struct {
	User       string `yaml:"user" env:"PG_USER" env-default:"postgres"`
	Password   string `yaml:"password" env:"PG_PASSWORD" env-default:"password"`
	Host       string `yaml:"host" env:"PG_HOST" env-default:"0.0.0.0"`
	Port       string `yaml:"port" env:"PG_PORT" env-default:"5432"`
	Database   string `yaml:"database" env:"PG_DATABASE" env-default:"postgres"`
	SSLMode    string `yaml:"ssl_mode" env:"PG_SSL" env-default:"disable"`
}

type JWT struct {
	AccessSignKey  string        `yaml:"access_sign_key" env:"ACCESS_JWT_KEY" env-default:"access_secret"`
	AccessTokenTTL time.Duration `yaml:"access_token_ttl" env:"ACCESS_JWT_TTL" env-default:"60m"`
	RefreshSignKey  string        `yaml:"refresh_sign_key" env:"REFRESH_JWT_KEY" env-default:"refresh_secret"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env:"REFRESH_JWT_TTL" env-default:"168h"`
}

func New() (*Config, error) {
	path := fetchConfigPath()

	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func fetchConfigPath() string {
	var path string

	if path = os.Getenv("CONFIG_PATH"); path == "" {
		path = defaultConfigFile
	}

	return path
}
