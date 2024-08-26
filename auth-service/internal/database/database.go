package database

import (
	"database/sql"
	"fmt"
	"go/mod/config"
	"log"
	"time"

	_ "github.com/lib/pq" // драйвер PostgreSQL
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// InitDB инициализирует подключение к базе данных
func InitDB(cfg *config.Config) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Проверяем подключение к базе
	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to ping the database: %v\n", err)
	}

	log.Println("Successfully connected to the database!")
}

// SaveRefreshToken сохраняет хэшированный Refresh токен для пользователя
func SaveRefreshToken(userID, hashedToken, ipAddress string) error {
	query := `INSERT INTO refresh_tokens (token_hash, user_id, ip_address) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, hashedToken, userID, ipAddress)
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %v", err)
	}
	return nil
}

// CheckRefreshToken проверяет корректность Refresh токена и IP-адреса
func CheckRefreshToken(userID, refreshToken, ipAddress string) (bool, error) {
	var storedHash, storedIP string
	query := `SELECT token_hash, ip_address FROM refresh_tokens WHERE user_id = $1`
	err := db.QueryRow(query, userID).Scan(&storedHash, &storedIP)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("refresh token not found")
		}
		return false, err
	}

	// Проверка совпадения токена
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(refreshToken))
	if err != nil {
		return false, fmt.Errorf("invalid refresh token")
	}

	// Проверка изменения IP-адреса
	if storedIP != ipAddress {
		log.Println("IP address has changed. Sending warning email...")
		// Здесь можно добавить логику для отправки письма
	}

	return true, nil
}

// UpdateRefreshToken обновляет Refresh токен в базе данных
func UpdateRefreshToken(userID, hashedToken, ipAddress string) error {
	query := `UPDATE refresh_tokens SET token_hash = $1, ip_address = $2, created_at = $3 WHERE user_id = $4`
	_, err := db.Exec(query, hashedToken, ipAddress, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update refresh token: %v", err)
	}
	return nil
}
