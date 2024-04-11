package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) getAuthRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("POST /login/", h.login)
	return router
}
