package handler

import "net/http"

func GetRouter() *http.ServeMux {
	router := http.NewServeMux()

	userRouter := getUserRouter()
	router.Handle("/users/", userRouter)
	return router
}
