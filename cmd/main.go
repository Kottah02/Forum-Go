package main

import (
	"log"
	"net/http"
	"website/internal/config"
	"website/internal/controllers"
	"website/internal/routes"
)

func main() {
	// Initialisation de la base de données
	config.InitDB()
	defer config.DB.Close()

	// Initialisation des contrôleurs
	authController := controllers.NewAuthController()
	userController := controllers.NewUserController()
	postController := controllers.NewPostController()

	// Configuration des routes
	routes.SetupRoutes(authController, userController, postController)

	// Démarrage du serveur
	log.Println("Serveur démarré sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Erreur de démarrage du serveur: %v", err)
	}
}
