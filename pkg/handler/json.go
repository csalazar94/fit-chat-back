package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func errorResponse(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	jsonResponse(w, code, errorResponse{
		Error: msg,
	})
}

func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error al serializar la respuesta: %v", err)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
