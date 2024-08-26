package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// HashToken генерирует bcrypt хэш для заданного токена
func HashToken(token string) (string, error) {
	// Генерация хэша с использованием bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CompareTokens сравнивает хэшированный токен с оригинальным токеном
func CompareTokens(token, hashedToken string) bool {
	// Проверка совпадения оригинального токена и его хэша
	err := bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token))
	if err != nil {
		log.Println("Token comparison failed:", err)
		return false
	}
	return true
}
