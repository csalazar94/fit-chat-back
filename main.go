package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/csalazar94/fit-chat-back/internal/db"
	"github.com/csalazar94/fit-chat-back/pkg/config"
	"github.com/csalazar94/fit-chat-back/pkg/handler"
	"github.com/csalazar94/fit-chat-back/pkg/service"
	_ "github.com/lib/pq"
	"github.com/sashabaranov/go-openai"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.NewConfig()

	conn, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatal(err)
	}
	queries := db.New(conn)

	openaiClient := openai.NewClient(cfg.OpenaiApiKey)

	services := service.NewServices(queries, openaiClient)

	v1Handler := handler.NewHandler(services)

	router := http.NewServeMux()
	router.Handle("/v1/", http.StripPrefix("/v1", v1Handler.Router))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	errc := make(chan error, 1)
	go func() {
		fmt.Printf("Servidor escuchando en el puerto %s\n", cfg.Port)
		errc <- server.ListenAndServe()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return server.Shutdown(ctx)
}
