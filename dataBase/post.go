package dataBase

import (
	"forum/structures"
	"fmt"
)

// CreatePost inserts a new post into the database
func CreatePost(userID int, title, content string) (int, error) {
	// Insert post into the database
	result, err := DB.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, title, content)
	if err != nil {
		return 0, fmt.Errorf("failed to create post: %v", err)
	}

	// Get the inserted post's ID
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	return int(postID), nil
}

// GetPosts retrieves all posts
func GetPosts() ([]structures.Post, error) {
	rows, err := DB.Query("SELECT id, user_id, title, content, created_at FROM posts")
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %v", err)
	}
	defer rows.Close()

	var posts []structures.Post
	for rows.Next() {
		var post structures.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}
