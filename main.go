package main

import (
	"fmt"
	"net/http"

	"github.com/csalazar94/fit-chat-back/pkg/handler"
	"github.com/csalazar94/fit-chat-back/pkg/service"
	_ "github.com/lib/pq"
)

func main() {
	config := loadConfig()

	services := service.NewService(config.dbQueries)
	v1Handler := handler.NewHandler(services)
	v1Router := v1Handler.GetRouter()
	router := http.NewServeMux()
	router.Handle("/v1/", http.StripPrefix("/v1", v1Router))
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.port),
		Handler: router,
	}

	fmt.Printf("Servidor escuchando en el puerto %s\n", config.port)
	server.ListenAndServe()
}
