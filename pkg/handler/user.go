package handler

import "net/http"

func getUser(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, 200, struct {
		Res string
	}{
		"hola",
	})
}

func getUserRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /", getUser)
	return router
}
