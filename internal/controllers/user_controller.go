package controllers

import (
	"net/http"
	"website/internal/middleware"
	"website/internal/models"
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
	user, err := models.GetUserByUsername(username)
	if err != nil {
		// If user is not found in DB, clear session and redirect to login
		middleware.ClearAuthSession(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postCount, err := models.GetPostCountByUserID(user.ID)
	if err != nil {
		// Handle error, but don't crash the app for a count
		postCount = 0 // Default to 0 if error
	}

	commentCount, err := models.GetCommentCountByUserID(user.ID)
	if err != nil {
		// Handle error, but don't crash the app for a count
		commentCount = 0 // Default to 0 if error
	}

	recentPosts, err := models.GetRecentPostsByUserID(user.ID, 5) // Get up to 5 recent posts
	if err != nil {
		// Handle error, but don't crash the app
		recentPosts = []models.Post{} // Default to empty slice if error
	}

	recentComments, err := models.GetRecentCommentsByUserID(user.ID, 5) // Get up to 5 recent comments
	if err != nil {
		// Handle error, but don't crash the app
		recentComments = []*models.Comment{} // Default to empty slice if error
	}

	data := map[string]interface{}{
		"User":           user,
		"PostCount":      postCount,
		"CommentCount":   commentCount,
		"RecentPosts":    recentPosts,
		"RecentComments": recentComments,
	}

	renderTemplate(w, r, "user/profile", data)
}

func (c *UserController) Index(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "index", nil)
}
