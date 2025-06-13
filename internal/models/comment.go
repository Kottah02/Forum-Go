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
