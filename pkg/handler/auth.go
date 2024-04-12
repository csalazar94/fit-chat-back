package handler

import (
	"encoding/json"
	"net/http"

	"github.com/csalazar94/fit-chat-back/pkg/service"
)

type AuthHandler struct {
	services *service.Services
}

func NewAuthRouter(services *service.Services) *http.ServeMux {
	router := http.NewServeMux()
	authHandler := &AuthHandler{services}
	router.HandleFunc("POST /login/", authHandler.login)
	return router
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	type bodySchema struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	body := bodySchema{}
	err := decoder.Decode(&body)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Error al decodificar el cuerpo de la petición")
		return
	}
	ok, err := h.services.AuthService.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error al autenticar el usuario")
		return
	}
	jsonResponse(w, http.StatusOK, ok)
}
