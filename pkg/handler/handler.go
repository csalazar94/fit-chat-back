package handler

import (
	"net/http"

	"github.com/csalazar94/fit-chat-back/pkg/service"
)

type Handler struct {
	Router *http.ServeMux
}

func NewHandler(services *service.Service) *Handler {
	router := http.NewServeMux()

	userRouter := NewUserRouter(services)
	router.Handle("/users/", http.StripPrefix("/users", userRouter))

	authRouter := NewAuthRouter(services)
	router.Handle("/auth/", http.StripPrefix("/auth", authRouter))

	return &Handler{
		Router: router,
	}
}
