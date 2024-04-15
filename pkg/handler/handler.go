package handler

import (
	"net/http"

	"github.com/csalazar94/fit-chat-back/pkg/service"
)

type Handler struct {
	Router http.Handler
}

func NewHandler(services *service.Services) *Handler {
	router := http.NewServeMux()

	wsHandler := NewWsHandler(services)
	router.HandleFunc("/ws/", wsHandler.subscribe)

	userRouter := NewUserRouter(services)
	router.Handle("/users/", http.StripPrefix("/users", userRouter))

	authRouter := NewAuthRouter(services)
	router.Handle("/auth/", http.StripPrefix("/auth", authRouter))

	messageRouter := NewMessageRouter(services)
	router.Handle("/messages/", http.StripPrefix("/messages", messageRouter))

	return &Handler{
		Router: LogRequestMiddleware(router),
	}
}
