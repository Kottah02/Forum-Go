package routes

import (
	"net/http"
	"strings"
	"website/internal/controllers"
	"website/internal/middleware"
)

func SetupRoutes(authController *controllers.AuthController, userController *controllers.UserController, postController *controllers.PostController) {
	// Routes statiques
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes publiques
	http.HandleFunc("/", userController.Index)
	http.HandleFunc("/login", authController.Login)
	http.HandleFunc("/register", authController.Register)
	http.HandleFunc("/logout", authController.Logout)

	// Routes protégées
	http.HandleFunc("/profile", middleware.RequireAuth(userController.Profile))
	http.HandleFunc("/posts/create", middleware.RequireAuth(postController.Create))
	http.HandleFunc("/posts", postController.List)
	http.HandleFunc("/posts/consulter", postController.Show)
	http.HandleFunc("/posts/add-comment", middleware.RequireAuth(postController.AddComment))
	http.HandleFunc("/posts/delete", middleware.RequireAuth(postController.Delete))
	http.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/react") {
			middleware.RequireAuth(postController.React)(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}
