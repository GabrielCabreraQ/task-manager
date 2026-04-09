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
	URI    string
	DBName string
}

// Load carga la configuración desde variables de entorno y devuelve el objeto Config.
func Load() (*Config, error) {
	// Intenta cargar el archivo .env si existe.
	if err := godotenv.Load(); err != nil {
		log.Println("No se encuentra .env")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", ""),
		},
		MongoDB: MongoConfig{
			URI:    buildMongoURI(),
			DBName: getEnv("DB_NAME", ""),
		},
	}, nil
}

// buildMongoURI construye la cadena de conexión de MongoDB usando variables de entorno.
// No se encuentra implementada el DB_USER y DB_PASSWORD, por lo que se asume que la conexión es sin autenticación.
func buildMongoURI() string {
	user := getEnv("DB_USER", "")
	password := getEnv("DB_PASSWORD", "")
	host := getEnv("DB_URL", "")
	dbName := getEnv("DB_NAME", "")

	// Si se proporcionan usuario y contraseña, utiliza autenticación.
	if user != "" && password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=admin", user, password, host, dbName)
	}

	// Si no hay credenciales, crea una URI sin autenticación para desarrollo local.
	return fmt.Sprintf("mongodb://%s/%s", host, dbName)
}

// getEnv obtiene el valor de la variable de entorno o devuelve un valor por defecto.
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return fallback
}
