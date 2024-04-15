package handler

import (
	"net/http"

	"github.com/csalazar94/fit-chat-back/pkg/service"
	"github.com/google/uuid"
)

type MessageHandler struct {
	services *service.Services
}

func NewMessageRouter(services *service.Services) *http.ServeMux {
	router := http.NewServeMux()
	messageHandler := &MessageHandler{services}
	router.HandleFunc("GET /", messageHandler.getAllByChatId)
	return router
}

func (h *MessageHandler) getAllByChatId(w http.ResponseWriter, r *http.Request) {
	chatId := r.URL.Query().Get("chat_id")
	chatIdUUID, err := uuid.Parse(chatId)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "El chat_id no es válido")
		return
	}
	messages, err := h.services.MessageService.GetAllByChatID(
		r.Context(),
		service.GetAllByChatIDParams{ChatID: chatIdUUID},
	)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error al obtener mensajes")
		return
	}
	jsonResponse(w, http.StatusOK, messages)
}
