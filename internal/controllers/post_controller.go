package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"website/internal/middleware"
	"website/internal/models"
)

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

func (c *PostController) List(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to list posts.")
	posts, err := models.GetAllPosts()
	if err != nil {
		log.Printf("Error getting all posts: %v", err)
		http.Error(w, "Erreur lors de la récupération des posts", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully retrieved %d posts.", len(posts))

	if middleware.IsAuthenticated(r) {
		log.Println("User is authenticated. Retrieving user info...")
		username, _ := middleware.GetUserInfo(r)
		user, err := models.GetUserByUsername(username)
		if err != nil {
			log.Printf("Error getting user by username %s: %v", username, err)
			http.Error(w, "Erreur lors de la récupération de l'utilisateur", http.StatusInternalServerError)
			return
		}
		log.Printf("Successfully retrieved user: %s (ID: %d)", user.Username, user.ID)

		if err := models.GetUserReactions(user.ID, posts); err != nil {
			log.Printf("Error getting user reactions for user ID %d: %v", user.ID, err)
			http.Error(w, "Erreur lors de la récupération des réactions", http.StatusInternalServerError)
			return
		}
		log.Println("Successfully retrieved user reactions.")
	} else {
		log.Println("User is not authenticated.")
	}

	// Récupérer tous les tags disponibles
	tags, err := models.GetAllTags()
	if err != nil {
		log.Printf("Error getting tags: %v", err)
		http.Error(w, "Erreur lors de la récupération des tags", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Posts": posts,
		"Tags":  tags,
	}
	log.Println("Rendering template posts/list...")
	renderTemplate(w, r, "posts/list", data)
	log.Println("Template rendering completed.")
}

func (c *PostController) Create(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	if r.Method != "POST" {
		// Récupérer les tags pour le formulaire
		tags, err := models.GetAllTags()
		if err != nil {
			http.Error(w, "Erreur lors de la récupération des tags", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Tags": tags,
		}
		renderTemplate(w, r, "posts/create", data)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	tagIDs := r.Form["tags"] // Récupère tous les tags sélectionnés

	if title == "" || content == "" {
		http.Error(w, "Le titre et le contenu ne peuvent pas être vides", http.StatusBadRequest)
		return
	}

	// Convertir les tagIDs en []int
	var tagIDsInt []int
	for _, tagID := range tagIDs {
		id, err := strconv.Atoi(tagID)
		if err != nil {
			http.Error(w, "ID de tag invalide", http.StatusBadRequest)
			return
		}
		tagIDsInt = append(tagIDsInt, id)
	}

	username, _ := middleware.GetUserInfo(r)
	user, err := models.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération de l'utilisateur", http.StatusInternalServerError)
		return
	}

	if err := models.CreatePost(user.ID, title, content, tagIDsInt); err != nil {
		http.Error(w, "Erreur lors de la création du post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

func (c *PostController) React(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	postID := strings.TrimPrefix(r.URL.Path, "/posts/")
	postID = strings.TrimSuffix(postID, "/react")
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "ID de post invalide", http.StatusBadRequest)
		return
	}

	reactionType := r.FormValue("reaction_type")
	if reactionType != "like" && reactionType != "dislike" {
		http.Error(w, "Type de réaction invalide", http.StatusBadRequest)
		return
	}

	username, _ := middleware.GetUserInfo(r)
	user, err := models.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération de l'utilisateur", http.StatusInternalServerError)
		return
	}

	if err := models.ToggleReaction(postIDInt, user.ID, reactionType); err != nil {
		http.Error(w, "Erreur lors de la mise à jour de la réaction", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/posts", http.StatusSeeOther)
}

func (c *PostController) Show(w http.ResponseWriter, r *http.Request) {
	log.Printf("Attempting to show post. URL: %s", r.URL.String())
	idStr := r.URL.Query().Get("id")
	log.Printf("Extracted ID string: %s", idStr)
	if idStr == "" {
		http.NotFound(w, r)
		log.Println("ID string is empty.")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Error converting ID to int: %v", err)
		http.NotFound(w, r)
		return
	}
	log.Printf("Converted ID: %d", id)

	post, err := models.GetPostByID(id)
	if err != nil {
		log.Printf("Error getting post by ID %d: %v", id, err)
		http.Error(w, "Erreur lors de la récupération du post", http.StatusInternalServerError)
		return
	}

	if post == nil {
		log.Printf("Post with ID %d not found.", id)
		http.NotFound(w, r)
		return
	}
	log.Printf("Successfully retrieved post: %v", post.Title)

	// Récupérer les commentaires du post
	comments, err := models.GetCommentsByPostID(id)
	fmt.Println(comments, err)
	if err != nil {
		log.Printf("Error getting comments for post ID %d: %v", id, err)
		http.Error(w, "Erreur lors de la récupération des commentaires", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully retrieved %d comments for post ID %d.", len(comments), id)

	data := map[string]interface{}{
		"Post":     post,
		"Comments": comments,
	}
	log.Println("Rendering template posts/show...")
	renderTemplate(w, r, "posts/show", data)
	log.Println("Template rendering completed.")
}

func (c *PostController) AddComment(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to add comment.")
	if !middleware.IsAuthenticated(r) {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		log.Println("AddComment: User not authenticated.")
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		log.Println("AddComment: Invalid HTTP method.")
		return
	}

	idStr := r.URL.Query().Get("id")
	log.Printf("AddComment: Extracted Post ID string: %s", idStr)
	if idStr == "" {
		http.NotFound(w, r)
		log.Println("AddComment: Post ID string is empty.")
		return
	}

	postID, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		log.Printf("AddComment: Error converting Post ID to int: %v", err)
		return
	}
	log.Printf("AddComment: Converted Post ID: %d", postID)

	content := strings.TrimSpace(r.FormValue("content"))
	log.Printf("AddComment: Comment content: %s", content)
	if content == "" {
		http.Error(w, "Le contenu du commentaire ne peut pas être vide", http.StatusBadRequest)
		log.Println("AddComment: Comment content is empty.")
		return
	}

	username, _ := middleware.GetUserInfo(r)
	log.Printf("AddComment: Authenticated username: %s", username)
	user, err := models.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération de l'utilisateur", http.StatusInternalServerError)
		log.Printf("AddComment: Error getting user by username %s: %v", username, err)
		return
	}
	log.Printf("AddComment: Retrieved User ID: %d", user.ID)

	if err := models.CreateComment(postID, user.ID, content); err != nil {
		http.Error(w, "Erreur lors de la création du commentaire", http.StatusInternalServerError)
		log.Printf("AddComment: Error creating comment: %v", err)
		return
	}
	log.Println("AddComment: Comment created successfully.")

	http.Redirect(w, r, "/posts/consulter?id="+idStr, http.StatusSeeOther)
	log.Printf("AddComment: Redirecting to /posts/consulter?id=%s", idStr)
}
