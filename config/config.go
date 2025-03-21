package config

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Config holds the credentials to establish a connection with the db
type Config struct {
	DbName   string `env:"DB_NAME,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT,required"`
	SslMode  string `env:"DB_SSLMODE,required"`
}

// Loads the env and returns a db configuration
func loadConfig() (*Config, error) {
	//load env variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env: %s", err)
	}

	// Parse the env variable to a config struct
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Error parsing env to Config struct: %s", err)
		return nil, err
	}

	return cfg, nil
}

// Establishes a connection to a postgres database
func SetupConnection() *sql.DB {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Db config didn't load, Error: %s", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName, cfg.SslMode)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatalf("Database connection failed, details: %s", err)
	}

	return db
}
