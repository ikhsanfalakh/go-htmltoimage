package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Environment string
	AppURL      string
	AppPort     string
	DBUser      string
	DBPassword  string
	DBHost      string
	DBPort      string
	DBName      string
}

var AppEnv *Env

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppEnv = &Env{
		Environment: GetEnv("ENVIRONMENT", "development"),
		AppURL:      GetEnv("APP_URL", "http://localhost"),
		AppPort:     GetEnv("APP_PORT", "8080"),

		DBUser:     GetEnv("DB_USER", "root"),
		DBPassword: GetEnv("DB_PASSWORD", ""),
		DBHost:     GetEnv("DB_HOST", "localhost"),
		DBPort:     GetEnv("DB_PORT", "3306"),
		DBName:     GetEnv("DB_NAME", "db_name"),
	}
}
