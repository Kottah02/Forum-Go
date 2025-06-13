package controllers

import (
	"net/http"
	"strings"
	"website/internal/middleware"
	"website/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		valid, err := models.ValidateUser(username, password)
		if err != nil || !valid {
			renderTemplate(w, r, "auth/login", map[string]interface{}{
				"Error": "Nom d'utilisateur ou mot de passe incorrect",
			})
			return
		}

		if err := middleware.SetAuthSession(w, r, username); err != nil {
			http.Error(w, "Erreur lors de la création de la session", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	renderTemplate(w, r, "auth/login", nil)
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	if middleware.IsAuthenticated(r) {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		username := strings.TrimSpace(r.FormValue("username"))
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm-password")

		errorMsg, isValid := validateRegistration(username, password, confirmPassword)
		if !isValid {
			renderTemplate(w, r, "auth/register", map[string]interface{}{
				"Error": errorMsg,
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			renderTemplate(w, r, "auth/register", map[string]interface{}{
				"Error": "Une erreur est survenue lors de l'inscription",
			})
			return
		}

		if err := models.CreateUser(username, string(hashedPassword)); err != nil {
			renderTemplate(w, r, "auth/register", map[string]interface{}{
				"Error": "Une erreur est survenue lors de l'inscription",
			})
			return
		}

		renderTemplate(w, r, "auth/register", map[string]interface{}{
			"Success": "Inscription réussie ! Vous pouvez maintenant vous connecter.",
		})
		return
	}

	renderTemplate(w, r, "auth/register", nil)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	if err := middleware.ClearAuthSession(w, r); err != nil {
		http.Error(w, "Erreur lors de la déconnexion", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func validateRegistration(username, password, confirmPassword string) (string, bool) {
	if len(username) < 3 || len(username) > 50 {
		return "Le nom d'utilisateur doit contenir entre 3 et 50 caractères", false
	}

	if len(password) < 6 {
		return "Le mot de passe doit contenir au moins 6 caractères", false
	}

	if password != confirmPassword {
		return "Les mots de passe ne correspondent pas", false
	}

	exists, err := models.UserExists(username)
	if err != nil {
		return "Une erreur est survenue", false
	}
	if exists {
		return "Ce nom d'utilisateur est déjà pris", false
	}

	return "", true
}
