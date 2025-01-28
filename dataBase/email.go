package dataBase

import (
	"fmt"
	"forum/structures"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser inserts a new user into the database
func InsertUser(email, username, password string) error {
	// Hash the password
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	fmt.Println("Hashed password:", string(hashed_password))

	// Insert user into database
	_, err = DB.Exec("INSERT INTO users (email, username, hashed_password) VALUES (?, ?, ?)", email, username, string(hashed_password))
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	return nil
}

// GetUserByEmail retrieves a user by their email
func GetUserByEmail(email string) (*structures.User, error) {

	fmt.Println("Getting user with email:", email)

	row := DB.QueryRow("SELECT id, email, username, hashed_password FROM users WHERE email = ?", email)

	var user structures.User

	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.HashedPassword)

	fmt.Println("user.ID: ", user.ID)
	fmt.Println("user.Email: ", user.Email)
	fmt.Println("user.Username: ", user.Username)
	fmt.Println("user.HashedPassword: ", user.HashedPassword)

	if err != nil {
		fmt.Println("email not found")
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}

	return &user, nil
}
