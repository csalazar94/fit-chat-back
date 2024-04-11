package handler

import (
	"net/http"

	"github.com/csalazar94/fit-chat-back/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services}
}

func (h Handler) GetRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("/users/", http.StripPrefix("/users", h.getUserRouter()))
	router.Handle("/auth/", http.StripPrefix("/auth", h.getAuthRouter()))

	return router
}