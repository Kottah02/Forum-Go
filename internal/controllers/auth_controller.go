package controllers

import (
	"net/http"
	"strings"
	"website/internal/middleware"
	"website/internal/models"
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

	if r.Method != "POST" {
		renderTemplate(w, r, "auth/register", nil)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")

	// Validation des champs
	if username == "" || email == "" || password == "" {
		renderTemplate(w, r, "auth/register", map[string]interface{}{
			"Error": "Tous les champs sont obligatoires",
		})
		return
	}

	// Validation basique de l'email
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		renderTemplate(w, r, "auth/register", map[string]interface{}{
			"Error": "L'adresse email n'est pas valide",
		})
		return
	}

	if password != confirmPassword {
		renderTemplate(w, r, "auth/register", map[string]interface{}{
			"Error": "Les mots de passe ne correspondent pas",
		})
		return
	}

	// Vérifier si l'utilisateur existe déjà
	exists, err := models.UserExists(username)
	if err != nil {
		http.Error(w, "Erreur lors de la vérification de l'utilisateur", http.StatusInternalServerError)
		return
	}
	if exists {
		renderTemplate(w, r, "auth/register", map[string]interface{}{
			"Error": "Ce nom d'utilisateur est déjà pris",
		})
		return
	}

	// Vérifier si l'email existe déjà
	emailExists, err := models.EmailExists(email)
	if err != nil {
		http.Error(w, "Erreur lors de la vérification de l'email", http.StatusInternalServerError)
		return
	}
	if emailExists {
		renderTemplate(w, r, "auth/register", map[string]interface{}{
			"Error": "Cette adresse email est déjà utilisée",
		})
		return
	}

	// Créer l'utilisateur
	err = models.CreateUser(username, email, password)
	if err != nil {
		var errorMessage string
		switch err {
		case models.ErrPasswordTooShort:
			errorMessage = "Le mot de passe doit contenir au moins 12 caractères"
		case models.ErrPasswordNoUpper:
			errorMessage = "Le mot de passe doit contenir au moins une majuscule"
		case models.ErrPasswordNoLower:
			errorMessage = "Le mot de passe doit contenir au moins une minuscule"
		case models.ErrPasswordNoNumber:
			errorMessage = "Le mot de passe doit contenir au moins un chiffre"
		case models.ErrPasswordNoSpecial:
			errorMessage = "Le mot de passe doit contenir au moins un caractère spécial"
		default:
			errorMessage = "Erreur lors de la création du compte"
		}
		renderTemplate(w, r, "auth/register", map[string]interface{}{
			"Error": errorMessage,
		})
		return
	}

	// Rediriger vers la page de connexion
	http.Redirect(w, r, "/login", http.StatusSeeOther)
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
