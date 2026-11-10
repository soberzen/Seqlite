package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv   string
	HTTP     HTTPConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type HTTPConfig struct {
	Port int
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Schema   string
	SSLMode  string
}

type RedisConfig struct {
	Host       string
	Port       int
	Password   string
	DB         int
	TLSEnabled bool
}

type JWTConfig struct {
	Issuer     string
	Audience   string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func LoadConfig() (*Config, error) {
	port, err := envPort("HTTP_PORT", 8080)
	if err != nil {
		return nil, err
	}
	postgresPort, err := envPort("POSTGRES_PORT", 5432)
	if err != nil {
		return nil, err
	}
	redisPort, err := envPort("REDIS_PORT", 6379)
	if err != nil {
		return nil, err
	}
	redisDB, err := envInt("REDIS_DB", 0)
	if err != nil {
		return nil, err
	}
	if redisDB < 0 {
		return nil, fmt.Errorf("REDIS_DB must not be negative")
	}
	redisTLSEnabled, err := envBool("REDIS_TLS_ENABLED", false)
	if err != nil {
		return nil, err
	}

	for _, name := range []string{
		"POSTGRES_USER",
		"POSTGRES_NAME",
		"JWT_ISSUER",
		"JWT_AUDIENCE",
	} {
		if strings.TrimSpace(os.Getenv(name)) == "" {
			return nil, fmt.Errorf("%s is required", name)
		}
	}

	appEnv := envString("APP_ENV", "development")

	sslMode := "verify-full"
	if appEnv == "development" {
		sslMode = "disable"
	}

	cfg := &Config{
		AppEnv: appEnv,
		HTTP: HTTPConfig{
			Port: port,
		},
		Postgres: PostgresConfig{
			Host:     envString("POSTGRES_HOST", "localhost"),
			Port:     postgresPort,
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Name:     os.Getenv("POSTGRES_NAME"),
			Schema:   envString("POSTGRES_SCHEMA", "public"),
			SSLMode:  sslMode,
		},
		Redis: RedisConfig{
			Host:       envString("REDIS_HOST", "localhost"),
			Port:       redisPort,
			Password:   os.Getenv("REDIS_PASSWORD"),
			DB:         redisDB,
			TLSEnabled: redisTLSEnabled,
		},
		JWT: JWTConfig{
			Issuer:     os.Getenv("JWT_ISSUER"),
			Audience:   os.Getenv("JWT_AUDIENCE"),
			AccessTTL:  15 * time.Minute,
			RefreshTTL: 7 * 24 * time.Hour,
		},
	}
	return cfg, nil
}

func envString(name, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	return value
}

func envInt(name string, fallback int) (int, error) {
	value := os.Getenv(name)
	if value == "" {
		return fallback, nil
	}
	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", name)
	}
	return number, nil
}

func envPort(name string, fallback int) (int, error) {
	port, err := envInt(name, fallback)
	if err != nil {
		return 0, err
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("%s must be between 1 and 65535", name)
	}
	return port, nil
}

func envBool(name string, fallback bool) (bool, error) {
	value := os.Getenv(name)
	if value == "" {
		return fallback, nil
	}
	enabled, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean", name)
	}
	return enabled, nil
}
