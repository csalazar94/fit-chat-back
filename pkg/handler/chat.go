package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/csalazar94/fit-chat-back/pkg/service"
	"github.com/google/uuid"
)

type ChatHandler struct {
	services *service.Services
}

func NewChatRouter(services *service.Services) *http.ServeMux {
	router := http.NewServeMux()
	chatHandler := &ChatHandler{services}
	router.HandleFunc("POST /", chatHandler.create)
	router.HandleFunc("GET /", chatHandler.getAll)
	return router
}

func (h *ChatHandler) create(w http.ResponseWriter, r *http.Request) {
	type bodySchema struct {
		UserID uuid.UUID `json:"user_id"`
		Title  string    `json:"title"`
	}
	decoder := json.NewDecoder(r.Body)
	body := bodySchema{}
	err := decoder.Decode(&body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Error al decodificar el cuerpo de la petición")
		return
	}
	chat, err := h.services.ChatService.Create(r.Context(), service.CreateChatParams{
		ID:        uuid.New(),
		UserID:    body.UserID,
		Title:     body.Title,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error al crear chat")
		return
	}
	jsonResponse(w, http.StatusOK, chat)
}

func (h *ChatHandler) getAll(w http.ResponseWriter, r *http.Request) {
	chats, err := h.services.ChatService.GetAll(r.Context())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error al obtener chats")
		return
	}
	jsonResponse(w, http.StatusOK, chats)
}
