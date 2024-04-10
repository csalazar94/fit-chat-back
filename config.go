package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/csalazar94/fit-chat-back/internal/db"
	"github.com/joho/godotenv"
)

type config struct {
	port      string
	dbQueries *db.Queries
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

	return config{
		port:      port,
		dbQueries: queries,
	}
}
