package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type App struct {
	Env      string
	Port     string
	LogLevel string
}

type DB struct {
	Host         string
	Port         string
	Name         string
	User         string
	Password     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
}

type Cache struct {
	Driver     string
	DefaultTTL time.Duration
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}

type JWT struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type Config struct {
	App   App
	DB    DB
	Cache Cache
	Redis Redis
	JWT   JWT
}

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_PORT", "8080")
	v.SetDefault("APP_LOG_LEVEL", "info")

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_NAME", "escola_gestao")
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_SSL_MODE", "disable")
	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 5)

	v.SetDefault("CACHE_DRIVER", "memory")
	v.SetDefault("CACHE_DEFAULT_TTL", 300)

	v.SetDefault("REDIS_ADDR", "localhost:6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	v.SetDefault("JWT_SECRET_KEY", "change-me-in-production")
	v.SetDefault("JWT_ACCESS_TOKEN_TTL", 900)
	v.SetDefault("JWT_REFRESH_TOKEN_TTL", 604800)

	ttlSec := v.GetInt("CACHE_DEFAULT_TTL")
	accessTTLSec := v.GetInt("JWT_ACCESS_TOKEN_TTL")
	refreshTTLSec := v.GetInt("JWT_REFRESH_TOKEN_TTL")

	return &Config{
		App: App{
			Env:      v.GetString("APP_ENV"),
			Port:     v.GetString("APP_PORT"),
			LogLevel: v.GetString("APP_LOG_LEVEL"),
		},
		DB: DB{
			Host:         v.GetString("DB_HOST"),
			Port:         v.GetString("DB_PORT"),
			Name:         v.GetString("DB_NAME"),
			User:         v.GetString("DB_USER"),
			Password:     v.GetString("DB_PASSWORD"),
			SSLMode:      v.GetString("DB_SSL_MODE"),
			MaxOpenConns: v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns: v.GetInt("DB_MAX_IDLE_CONNS"),
		},
		Cache: Cache{
			Driver:     v.GetString("CACHE_DRIVER"),
			DefaultTTL: time.Duration(ttlSec) * time.Second,
		},
		Redis: Redis{
			Addr:     v.GetString("REDIS_ADDR"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
		},
		JWT: JWT{
			SecretKey:       v.GetString("JWT_SECRET_KEY"),
			AccessTokenTTL:  time.Duration(accessTTLSec) * time.Second,
			RefreshTokenTTL: time.Duration(refreshTTLSec) * time.Second,
		},
	}, nil
}
