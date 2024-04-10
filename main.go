package main

import (
	"fmt"
	"net/http"

	"github.com/csalazar94/fit-chat-back/pkg/handler"
)

func main() {
	config := loadConfig()

	v1Router := handler.GetRouter()
	router := http.NewServeMux()
	router.Handle("/v1/", http.StripPrefix("/v1", v1Router))
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: router,
	}

	fmt.Printf("Server listening on port %s", config.Port)
	server.ListenAndServe()
}
