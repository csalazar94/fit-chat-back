package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DbURL         string
	OpenaiApiKey  string
	GoogleKey     string
	GoogleSecret  string
	SessionSecret string
}

func NewConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error al cargar la variable de entorno 'PORT'")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Error al cargar la variable de entorno 'DB_URL'")
	}

	openaiApiKey := os.Getenv("OPENAI_API_KEY")
	if openaiApiKey == "" {
		log.Fatal("Error al cargar la variable de entorno 'OPENAI_API_KEY'")
	}

	googleKey := os.Getenv("GOOGLE_KEY")
	if googleKey == "" {
		log.Fatal("Error al cargar la variable de entorno 'GOOGLE_KEY'")
	}

	googleSecret := os.Getenv("GOOGLE_SECRET")
	if googleSecret == "" {
		log.Fatal("Error al cargar la variable de entorno 'GOOGLE_SECRET'")
	}

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Fatal("Error al cargar la variable de entorno 'SESSION_SECRET'")
	}

	return Config{
		Port:          port,
		DbURL:         dbURL,
		OpenaiApiKey:  openaiApiKey,
		GoogleKey:     googleKey,
		GoogleSecret:  googleSecret,
		SessionSecret: sessionSecret,
	}
}
