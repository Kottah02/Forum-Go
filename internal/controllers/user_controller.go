package controllers

import (
	"net/http"
	"website/internal/middleware"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (c *UserController) Profile(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username, _ := middleware.GetUserInfo(r)
	data := map[string]interface{}{
		"Username": username,
	}

	renderTemplate(w, r, "user/profile", data)
}

func (c *UserController) Index(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "index", nil)
}
