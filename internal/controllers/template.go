package controllers

import (
	"html/template"
	"log"
	"net/http"
	"website/internal/middleware"
)

func renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}

	// Ajouter les données de base pour toutes les pages
	if middleware.IsAuthenticated(r) {
		data["IsAuthenticated"] = true
		if username, ok := middleware.GetUserInfo(r); ok {
			data["Username"] = username
		}
	} else {
		data["IsAuthenticated"] = false
	}
	data["ActivePage"] = tmpl

	tmplFiles := []string{
		"templates/" + tmpl + ".html",
		"templates/nav.html",
	}

	t, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		log.Printf("Erreur lors du parsing des templates: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	if err = t.Execute(w, data); err != nil {
		log.Printf("Erreur lors de l'exécution du template: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
	}
}
