package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/csalazar94/fit-chat-back/pkg/service"
	"github.com/google/uuid"
)

type MessageHandler struct {
	services *service.Services
}

func NewMessageRouter(services *service.Services) *http.ServeMux {
	router := http.NewServeMux()
	messageHandler := &MessageHandler{services}
	router.HandleFunc("POST /", messageHandler.createMessage)
	return router
}

func (h *MessageHandler) createMessage(w http.ResponseWriter, r *http.Request) {
	type bodySchema struct {
		ChatID       uuid.UUID `json:"chat_id"`
		AuthorRoleID int32     `json:"author_role_id"`
		Content      string    `json:"content"`
	}
	decoder := json.NewDecoder(r.Body)
	body := bodySchema{}
	err := decoder.Decode(&body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Error al decodificar el cuerpo de la petición")
		return
	}
	message, err := h.services.MessageService.Create(r.Context(), service.CreateMessageParams{
		ID:           uuid.New(),
		ChatID:       body.ChatID,
		AuthorRoleID: body.AuthorRoleID,
		Content:      body.Content,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error al crear el mensaje")
		return
	}
	jsonResponse(w, http.StatusOK, message)
}
