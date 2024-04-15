package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/csalazar94/fit-chat-back/internal/db"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

type config struct {
	port         string
	dbQueries    *db.Queries
	openaiClient *openai.Client
}

func loadConfig() config {
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
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	queries := db.New(conn)

	openaiApiKey := os.Getenv("OPENAI_API_KEY")
	if openaiApiKey == "" {
		log.Fatal("Error al cargar la variable de entorno 'OPENAI_API_KEY'")
	}
	openaiClient := openai.NewClient(openaiApiKey)

	return config{
		port:         port,
		dbQueries:    queries,
		openaiClient: openaiClient,
	}
}
