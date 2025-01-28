package handlers

import (
	"fmt"
	"forum/dataBase"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	//"forum/structures"
)

// Login handler
func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "login.html", nil)
		return
	} else if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

	email := r.FormValue("email")
	password := r.FormValue("password")


	fmt.Println("login password: ", password)

	// Retrieve user from the database
	user, err := dataBase.GetUserByEmail(email) 
if err != nil || bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
		

	// Set session and redirect
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    fmt.Sprintf("%d", user.ID),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	})

	http.Redirect(w, r, "/home", http.StatusFound)
}

