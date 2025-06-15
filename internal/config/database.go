package config

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	DB    *sql.DB
	store = sessions.NewCookieStore([]byte("super-secret-key"))
)

func init() {
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 jours
		HttpOnly: true,
		Secure:   false, // Mettre à true en production avec HTTPS
		SameSite: http.SameSiteLaxMode,
	}
}

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
		email VARCHAR(191) UNIQUE NOT NULL,
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

	createTagsTableSQL := `
	CREATE TABLE IF NOT EXISTS tags (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	createPostTagsTableSQL := `
	CREATE TABLE IF NOT EXISTS post_tags (
		post_id INT NOT NULL,
		tag_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (post_id, tag_id),
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);`

	createCommentsTableSQL := `
	CREATE TABLE IF NOT EXISTS comments (
		id INT AUTO_INCREMENT PRIMARY KEY,
		post_id INT NOT NULL,
		user_id INT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	// Exécution des requêtes de création des tables
	tables := []string{
		createUsersTableSQL,
		createPostsTableSQL,
		createPostReactionsTableSQL,
		createTagsTableSQL,
		createPostTagsTableSQL,
		createCommentsTableSQL,
	}

	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			log.Printf("Erreur lors de la création d'une table: %v", err)
		}
	}

	// Migration pour ajouter le champ email s'il n'existe pas
	var columnExists bool
	err := DB.QueryRow(`
		SELECT COUNT(*) > 0
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_NAME = 'users' 
		AND COLUMN_NAME = 'email'
	`).Scan(&columnExists)

	if err != nil {
		log.Printf("Erreur lors de la vérification de la colonne email: %v", err)
	} else if !columnExists {
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN email VARCHAR(191) UNIQUE NOT NULL AFTER username")
		if err != nil {
			log.Printf("Erreur lors de l'ajout de la colonne email: %v", err)
		} else {
			log.Println("Colonne email ajoutée avec succès")
		}
	}

	// Insertion des tags par défaut
	insertDefaultTagsSQL := `
	INSERT IGNORE INTO tags (name) VALUES 
	('Arme'),
	('Skin'),
	('Gameplay'),
	('Ynov');`

	_, err = DB.Exec(insertDefaultTagsSQL)
	if err != nil {
		log.Printf("Erreur lors de l'insertion des tags par défaut: %v", err)
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
