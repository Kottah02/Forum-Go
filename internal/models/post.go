package models

import (
	"database/sql"
	"time"
	"website/internal/config"
)

type Post struct {
	ID           int
	Title        string
	Content      string
	Author       string
	CreatedAt    time.Time
	LikeCount    int
	DislikeCount int
	UserReaction string // "like", "dislike" ou ""
}

func GetAllPosts() ([]Post, error) {
	query := `
		SELECT 
			p.id, 
			p.title, 
			p.content, 
			u.username, 
			p.created_at,
			COALESCE((SELECT COUNT(*) FROM post_reactions WHERE post_id = p.id AND reaction_type = 'like'), 0) as like_count,
			COALESCE((SELECT COUNT(*) FROM post_reactions WHERE post_id = p.id AND reaction_type = 'dislike'), 0) as dislike_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var createdAt time.Time
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.Author,
			&createdAt,
			&post.LikeCount,
			&post.DislikeCount,
		)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = createdAt
		posts = append(posts, post)
	}

	return posts, rows.Err()
}

func CreatePost(userID int, title, content string) error {
	_, err := config.DB.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)",
		userID, title, content)
	return err
}

func GetUserReactions(userID int, posts []Post) error {
	for i := range posts {
		var reactionType string
		err := config.DB.QueryRow("SELECT reaction_type FROM post_reactions WHERE post_id = ? AND user_id = ?",
			posts[i].ID, userID).Scan(&reactionType)
		if err == nil {
			posts[i].UserReaction = reactionType
		} else if err != sql.ErrNoRows {
			return err
		}
	}
	return nil
}

func ToggleReaction(postID, userID int, reactionType string) error {
	var existingReaction string
	err := config.DB.QueryRow("SELECT reaction_type FROM post_reactions WHERE post_id = ? AND user_id = ?",
		postID, userID).Scan(&existingReaction)

	if err == sql.ErrNoRows {
		_, err = config.DB.Exec("INSERT INTO post_reactions (post_id, user_id, reaction_type) VALUES (?, ?, ?)",
			postID, userID, reactionType)
		return err
	} else if err != nil {
		return err
	}

	if existingReaction == reactionType {
		_, err = config.DB.Exec("DELETE FROM post_reactions WHERE post_id = ? AND user_id = ?",
			postID, userID)
	} else {
		_, err = config.DB.Exec("UPDATE post_reactions SET reaction_type = ? WHERE post_id = ? AND user_id = ?",
			reactionType, postID, userID)
	}
	return err
}

func GetPostByID(id int) (*Post, error) {
	row := config.DB.QueryRow("SELECT p.id, p.title, p.content, u.username, p.created_at, "+
		"COALESCE(SUM(CASE WHEN pr.reaction_type = 'like' THEN 1 ELSE 0 END), 0) AS likes, "+
		"COALESCE(SUM(CASE WHEN pr.reaction_type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislikes "+
		"FROM posts p "+
		"JOIN users u ON p.user_id = u.id "+
		"LEFT JOIN post_reactions pr ON p.id = pr.post_id "+
		"WHERE p.id = ? GROUP BY p.id", id)
	post := &Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CreatedAt, &post.LikeCount, &post.DislikeCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Post not found
		}
		return nil, err
	}
	return post, nil
}
