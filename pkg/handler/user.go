package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/csalazar94/fit-chat-back/pkg/service"
	"github.com/google/uuid"
)

type UserHandler struct {
	services *service.Services
}

func NewUserRouter(services *service.Services) *http.ServeMux {
	router := http.NewServeMux()
	userHandler := &UserHandler{services}
	router.HandleFunc("GET /", userHandler.getUsers)
	router.HandleFunc("POST /", userHandler.createUser)
	return router
}

func (h *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.services.UserService.GetAll(r.Context())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error al obtener los usuarios")
		return
	}
	if users != nil {
		jsonResponse(w, http.StatusOK, users)
		return
	}
	jsonResponse(w, http.StatusOK, []service.User{})
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	type bodySchema struct {
		FullName string `json:"full_name"`
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
	user, err := h.services.UserService.Create(r.Context(), service.CreateUserParams{
		ID:        uuid.New(),
		FullName:  body.FullName,
		Email:     body.Email,
		Password:  body.Password,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error al crear el usuario")
		return
	}
	jsonResponse(w, http.StatusOK, user)
}
