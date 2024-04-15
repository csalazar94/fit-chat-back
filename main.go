package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/csalazar94/fit-chat-back/pkg/handler"
	"github.com/csalazar94/fit-chat-back/pkg/service"
	_ "github.com/lib/pq"
)

func main() {
	config := loadConfig()

	err := run(config)
	if err != nil {
		log.Fatal(err)
	}
}

func run(config config) error {
	services := service.NewServices(config.dbQueries, config.openaiClient)
	v1Handler := handler.NewHandler(services)
	router := http.NewServeMux()

	router.Handle("/v1/", http.StripPrefix("/v1", v1Handler.Router))
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.port),
		Handler: router,
	}

	errc := make(chan error, 1)
	go func() {
		fmt.Printf("Servidor escuchando en el puerto %s\n", config.port)
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
