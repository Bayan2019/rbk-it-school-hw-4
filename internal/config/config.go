package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/time/rate"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Api      ApiConfig
}

type AppConfig struct {
	Port         string
	JwtSecret    string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type ApiConfig struct {
	Limiter   *rate.Limiter
	UserAgent string
}

func MustLoad() Config {
	cfg := Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			// 6. Безопасность
			// - JWT secret хранить в env
			JwtSecret:    getEnv("JWT_SECRET", "dev-secret-change-me"),
			ReadTimeout:  mustDuration("APP_READ_TIMEOUT", "5s"),
			WriteTimeout: mustDuration("APP_WRITE_TIMEOUT", "10s"),
			IdleTimeout:  mustDuration("APP_IDLE_TIMEOUT", "60s"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "users_db"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Api: ApiConfig{
			Limiter:   rate.NewLimiter(rate.Every(time.Second), 5),
			UserAgent: getEnv("USER_AGENT", "weather-api/1.0 (example@gmail.com)"),
		},
	}

	return cfg
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Host,
		c.Port,
		c.Name,
		c.User,
		c.Password,
		c.SSLMode,
	)
}

////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func mustDuration(key, fallback string) time.Duration {
	value := getEnv(key, fallback)
	d, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("invalid duration %s=%s: %v", key, value, err)
	}
	return d
}
