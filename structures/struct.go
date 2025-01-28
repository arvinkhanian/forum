package structures

import "time"

type Post struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
}

// User represents a user in the system
type User struct {
	ID       int
	Email    string
	Username string
	HashedPassword string // In practice, this should be a hashed password
}

var Users User 