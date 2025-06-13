package models

import (
	"time"
	"website/internal/config"
)

type Comment struct {
	ID        int
	PostID    int
	UserID    int
	Content   string
	Author    string
	CreatedAt time.Time
	PostTitle string
}

func CreateComment(postID, userID int, content string) error {
	_, err := config.DB.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
	return err
}

func GetCommentsByPostID(postID int) ([]*Comment, error) {
	rows, err := config.DB.Query("SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username FROM comments c JOIN users u ON c.user_id = u.id WHERE c.post_id = ? ORDER BY c.created_at ASC", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		comment := &Comment{}
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.Author)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func GetCommentCountByUserID(userID int) (int, error) {
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM comments WHERE user_id = ?", userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetRecentCommentsByUserID(userID int, limit int) ([]*Comment, error) {
	query := `
		SELECT 
			c.id, 
			c.post_id, 
			c.user_id, 
			c.content, 
			c.created_at, 
			u.username AS author_username, 
			p.title AS post_title
		FROM comments c
		JOIN users u ON c.user_id = u.id
		JOIN posts p ON c.post_id = p.id
		WHERE c.user_id = ?
		ORDER BY c.created_at DESC
		LIMIT ?
	`

	rows, err := config.DB.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		comment := &Comment{}
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.Author, &comment.PostTitle)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, rows.Err()
}
