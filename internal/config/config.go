package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server  ServerConfig
	MongoDB MongoConfig
}

type ServerConfig struct {
	Port string
}

type MongoConfig struct {
	URI string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encuentra .env")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		MongoDB: MongoConfig{
			URI: buildMongoURI(),
		},
	}, nil
}

// buildMongoURI mejorada con authSource
func buildMongoURI() string {
	//user := getEnv("DB_USER", "")
	//password := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "")
	host := getEnv("DB_URL", "")

	// Si hay usuario y contraseña → autenticación
	//if user != "" && password != "" {
	//	return fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=admin", user, password, host, dbName)
	//}

	// Sin autenticación
	return fmt.Sprintf("mongodb://%s/%s", host, dbName)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return fallback
}
