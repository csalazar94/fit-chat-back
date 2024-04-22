package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/csalazar94/fit-chat-back/pkg/service"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthHandler struct {
	services  *service.Services
	oauthConf *oauth2.Config
}

func NewAuthRouter(services *service.Services) *http.ServeMux {
	router := http.NewServeMux()
	authHandler := &AuthHandler{
		services: services, oauthConf: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_KEY"),
			ClientSecret: os.Getenv("GOOGLE_SECRET"),
			RedirectURL:  "http://localhost:3000/v1/auth/google/callback",
			Endpoint:     google.Endpoint,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		},
	}
	router.HandleFunc("GET /google/callback", authHandler.googleCallback)
	router.HandleFunc("GET /google/logout", authHandler.googleLogout)
	router.HandleFunc("GET /google/login", authHandler.googleLogin)
	return router
}

func (h *AuthHandler) googleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, err := h.oauthConf.Exchange(r.Context(), code)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error al intercambiar el código por un token")
	}
	client := h.oauthConf.Client(r.Context(), token)
	response, err := client.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token.AccessToken))
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error al obtener la información del usuario")
	}
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if response.StatusCode != 200 {
		errorResponse(w, http.StatusInternalServerError, "Error al obtener la información del usuario")
	}
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Error al obtener la información del usuario")
	}
	jsonResponse(w, 200, string(body))
}

func (h *AuthHandler) googleLogout(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 200, "logout")
}

func (h *AuthHandler) googleLogin(w http.ResponseWriter, r *http.Request) {
	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		panic("source of randomness unavailable: " + err.Error())
	}
	url := h.oauthConf.AuthCodeURL(base64.URLEncoding.EncodeToString(nonceBytes))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
