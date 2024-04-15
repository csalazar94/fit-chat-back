package handler

import (
	"net/http"

	"github.com/csalazar94/fit-chat-back/pkg/service"
)

type ChatHandler struct {
	services *service.Services
}

func NewChatRouter(services *service.Services) *http.ServeMux {
	router := http.NewServeMux()
	chatHandler := &ChatHandler{services}
	router.HandleFunc("GET /", chatHandler.getAll)
	return router
}

func (h *ChatHandler) getAll(w http.ResponseWriter, r *http.Request) {
	chats, err := h.services.ChatService.GetAll(r.Context())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error al obtener chats")
		return
	}
	jsonResponse(w, http.StatusOK, chats)
}
