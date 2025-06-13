package models

import (
	"database/sql"
	"time"
	"website/internal/config"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int
	Username  string
	CreatedAt time.Time
}

func GetUserByUsername(username string) (User, error) {
	var user User
	err := config.DB.QueryRow("SELECT id, username, created_at FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, err
	}
	return user, nil
}

func CreateUser(username, password string) error {
	_, err := config.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	return err
}

func ValidateUser(username, password string) (bool, error) {
	var hashedPassword string
	err := config.DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil, nil
}

func UserExists(username string) (bool, error) {
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	return exists, err
}
