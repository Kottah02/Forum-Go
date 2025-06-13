package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	DB    *sql.DB
	store = sessions.NewCookieStore([]byte("super-secret-key"))
)

func GetSessionStore() *sessions.CookieStore {
	return store
}

func InitDB() {
	var err error
	DB, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/website_db?parseTime=true")
	if err != nil {
		log.Fatalf("Erreur de connexion à la base de données: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Erreur de ping de la base de données: %v", err)
	}
	log.Println("Connexion à la base de données réussie!")

	createTables()
	createTestUser()
}

func createTables() {
	createUsersTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	createPostsTableSQL := `
	CREATE TABLE IF NOT EXISTS posts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	createPostReactionsTableSQL := `
	CREATE TABLE IF NOT EXISTS post_reactions (
		id INT AUTO_INCREMENT PRIMARY KEY,
		post_id INT NOT NULL,
		user_id INT NOT NULL,
		reaction_type ENUM('like', 'dislike') NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE KEY unique_user_post (post_id, user_id)
	);`

	_, err := DB.Exec(createUsersTableSQL)
	if err != nil {
		log.Printf("Erreur lors de la création de la table users: %v", err)
	}

	_, err = DB.Exec(createPostsTableSQL)
	if err != nil {
		log.Printf("Erreur lors de la création de la table posts: %v", err)
	}

	_, err = DB.Exec(createPostReactionsTableSQL)
	if err != nil {
		log.Printf("Erreur lors de la création de la table post_reactions: %v", err)
	}
}

func createTestUser() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", "test").Scan(&count)
	if err != nil {
		log.Printf("Erreur lors de la vérification de l'utilisateur de test: %v", err)
		return
	}

	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
		_, err = DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", "test", string(hashedPassword))
		if err != nil {
			log.Printf("Erreur lors de l'insertion de l'utilisateur de test: %v", err)
		} else {
			log.Println("Utilisateur de test créé avec succès")
		}
	} else {
		log.Println("L'utilisateur test existe déjà")
	}

	// Créer un post de test si la table est vide
	err = DB.QueryRow("SELECT COUNT(*) FROM posts").Scan(&count)
	if err != nil {
		log.Printf("Erreur lors de la vérification des posts existants: %v", err)
		return
	}

	if count == 0 {
		var testUserID int
		err = DB.QueryRow("SELECT id FROM users WHERE username = ?", "test").Scan(&testUserID)
		if err != nil {
			log.Printf("Erreur lors de la récupération de l'ID de l'utilisateur test: %v", err)
			return
		}

		_, err = DB.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)",
			testUserID, "Pourquoi Cyril est beau", "Pourquoi pas finalement")
		if err != nil {
			log.Printf("Erreur lors de l'insertion du post de test: %v", err)
		} else {
			log.Println("Post de test créé avec succès")
		}
	}
}
