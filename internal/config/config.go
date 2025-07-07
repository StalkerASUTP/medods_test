package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}
type Secret struct {
	SecretKey  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}
type HTTPServer struct {
	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration
}
type Config struct {
	DB
	Secret
	HTTPServer
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	cfg := &Config{
		DB: DB{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			DBName:   os.Getenv("DB_NAME"),
		},
		Secret: Secret{
			SecretKey:  os.Getenv("JWT_SECRET"),
			AccessTTL:  parseDuration(os.Getenv("ACCESS_TTL")),
			RefreshTTL: parseDuration(os.Getenv("REFRESH_TTL")),
		},
		HTTPServer: HTTPServer{
			Address:     os.Getenv("SERVER_ADDRESS"),
			Timeout:     parseDuration(os.Getenv("SERVER_TIMEOUT")),
			IdleTimeout: parseDuration(os.Getenv("SERVER_IDLE_TIMEOUT")),
		},
	}
	return cfg
}
func parseDuration(times string) time.Duration {
	parsedTime, err := time.ParseDuration(times)
	if err != nil {
		panic(fmt.Sprintf("Error parsing duration: %v", err))
	}
	return parsedTime
}
