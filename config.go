package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	Port string
}

func loadConfig() config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error loading port environment variable")
	}

	return config{
		Port: port,
	}
}
