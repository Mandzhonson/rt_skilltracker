package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath = ".env"

type Config struct {
	HTTP     HTTPConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Minio    MinioConfig
}

type HTTPConfig struct {
	Port            string        `env:"HTTP_PORT" env-default:"8080"`
	Host            string        `env:"HTTP_HOST" env-default:"0.0.0.0"`
	ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"5s"`
	WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"10s"`
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" env-required:"true"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Name     string `env:"POSTGRES_DB" env-required:"true"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" env-required:"true"`
	Port     string `env:"REDIS_PORT" env-default:"6379"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
}

type MinioConfig struct {
	User     string `env:"MINIO_ROOT_USER" env-required:"true"`
	Password string `env:"MINIO_ROOT_PASSWORD" env-required:"true"`
	Host     string `env:"MINIO_HOST" env-required:"true"`
	Port     string `env:"MINIO_API_PORT" env-default:"9000"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		if err := cleanenv.ReadEnv(cfg); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

func (cfg Config) PostgresDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",cfg.Postgres.User,cfg.Postgres.Password,cfg.Postgres.Host,cfg.Postgres.Port,cfg.Postgres.Name)
}
