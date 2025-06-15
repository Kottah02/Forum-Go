package models

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"time"
	"website/internal/config"
)

type User struct {
	ID        int
	Username  string
	CreatedAt time.Time
}

var (
	ErrPasswordTooShort  = errors.New("le mot de passe doit contenir au moins 12 caractères")
	ErrPasswordNoUpper   = errors.New("le mot de passe doit contenir au moins une majuscule")
	ErrPasswordNoLower   = errors.New("le mot de passe doit contenir au moins une minuscule")
	ErrPasswordNoNumber  = errors.New("le mot de passe doit contenir au moins un chiffre")
	ErrPasswordNoSpecial = errors.New("le mot de passe doit contenir au moins un caractère spécial")
)

func GetUserByUsername(username string) (User, error) {
	var user User
	err := config.DB.QueryRow("SELECT id, username, created_at FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("utilisateur non trouvé: %s", username)
		}
		return User{}, err
	}
	return user, nil
}

func validatePassword(password string) error {
	if len(password) < 12 {
		return ErrPasswordTooShort
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return ErrPasswordNoUpper
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return ErrPasswordNoLower
	}

	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return ErrPasswordNoNumber
	}

	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	if !hasSpecial {
		return ErrPasswordNoSpecial
	}

	return nil
}

func CreateUser(username, email, password string) error {
	// Valider le mot de passe
	if err := validatePassword(password); err != nil {
		return err
	}

	// Hacher le mot de passe avec SHA-512
	hasher := sha512.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	// Insérer l'utilisateur avec le mot de passe haché
	_, err := config.DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		username, email, hashedPassword)
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

	// Hacher le mot de passe fourni avec SHA-512
	hasher := sha512.New()
	hasher.Write([]byte(password))
	providedHash := hex.EncodeToString(hasher.Sum(nil))

	// Comparer les hash
	return providedHash == hashedPassword, nil
}

func UserExists(username string) (bool, error) {
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	return exists, err
}

func EmailExists(email string) (bool, error) {
	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
	return exists, err
}

func UpdatePassword(userID int, newPassword string) error {
	// Valider le nouveau mot de passe
	if err := validatePassword(newPassword); err != nil {
		return err
	}

	// Le reste du code de mise à jour du mot de passe...
	// ... existing code ...
	return nil
}
