package models

import (
	"database/sql"
	"fmt"
	"log"
)

// User представляет модель пользователя
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// GetUserByID получает пользователя по его идентификатору (GUID)
func GetUserByID(db *sql.DB, userID string) (*User, error) {
	var user User
	query := `SELECT id, email, name FROM users WHERE id = $1`
	err := db.QueryRow(query, userID).Scan(&user.ID, &user.Email, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %s not found", userID)
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser создает нового пользователя в базе данных
func CreateUser(db *sql.DB, user *User) error {
	query := `INSERT INTO users (id, email, name) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, user.ID, user.Email, user.Name)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

// UpdateUser обновляет информацию о пользователе в базе данных
func UpdateUser(db *sql.DB, user *User) error {
	query := `UPDATE users SET email = $1, name = $2 WHERE id = $3`
	_, err := db.Exec(query, user.Email, user.Name, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

// DeleteUser удаляет пользователя из базы данных по его идентификатору
func DeleteUser(db *sql.DB, userID string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}

// LogUserInfo выводит информацию о пользователе в лог
func LogUserInfo(user *User) {
	log.Printf("User Info: ID=%s, Email=%s, Name=%s\n", user.ID, user.Email, user.Name)
}
