package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"go/mod/internal/bcrypt"
	"go/mod/internal/database"
	"go/mod/internal/email"
	"go/mod/internal/jwt"
	"log"
	"net/http"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GenerateTokens обрабатывает запрос для генерации Access и Refresh токенов
func GenerateTokens(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	// Создаем Access и Refresh токены
	accessToken, err := jwt.GenerateAccessToken(userID, r.RemoteAddr)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken := uuid.New().String()

	// Хэшируем Refresh токен
	hashedRefreshToken, err := bcrypt.HashToken(refreshToken)
	if err != nil {
		http.Error(w, "Failed to hash refresh token", http.StatusInternalServerError)
		return
	}

	// Сохраняем Refresh токен в базе данных
	err = database.SaveRefreshToken(userID, hashedRefreshToken, r.RemoteAddr)
	if err != nil {
		http.Error(w, "Failed to save refresh token", http.StatusInternalServerError)
		return
	}

	// Возвращаем токены клиенту
	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RefreshToken обрабатывает запрос для обновления Access токена через Refresh токен
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID       string `json:"user_id"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.RefreshToken == "" {
		http.Error(w, "Missing user_id or refresh_token", http.StatusBadRequest)
		return
	}

	// Проверяем Refresh токен в базе данных
	ipAddress := r.RemoteAddr
	valid, err := database.CheckRefreshToken(req.UserID, req.RefreshToken, ipAddress)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if !valid {
		http.Error(w, "Refresh token is invalid or expired", http.StatusUnauthorized)
		return
	}

	// Генерируем новый Access токен
	newAccessToken, err := jwt.GenerateAccessToken(req.UserID, ipAddress)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Генерируем новый Refresh токен
	newRefreshToken := uuid.New().String()

	// Хэшируем новый Refresh токен
	hashedNewRefreshToken, err := bcrypt.HashToken(newRefreshToken)
	if err != nil {
		http.Error(w, "Failed to hash new refresh token", http.StatusInternalServerError)
		return
	}

	// Обновляем Refresh токен в базе данных
	err = database.UpdateRefreshToken(req.UserID, hashedNewRefreshToken, ipAddress)
	if err != nil {
		http.Error(w, "Failed to update refresh token", http.StatusInternalServerError)
		return
	}

	// Отправляем предупреждение, если IP-адрес изменился
	if r.RemoteAddr != ipAddress {
		log.Println("IP address has changed, sending warning email...")
		email.SendWarningEmail("user@example.com", ipAddress) // Mock email для примера
	}

	// Возвращаем новые токены клиенту
	response := TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
