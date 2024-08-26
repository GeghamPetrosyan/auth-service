package main

import (
	"go/mod/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Маршруты для работы с токенами
	r.HandleFunc("/auth/token", handlers.GenerateTokens).Methods("GET")
	r.HandleFunc("/auth/refresh", handlers.RefreshToken).Methods("POST")

	// Запуск сервера на порту 8080
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
