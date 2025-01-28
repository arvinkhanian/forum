package handlers

import (
	"fmt"
	"forum/dataBase"
	"html/template"
	"net/http"
	//"forum/structures"
)

var templates = template.Must(template.ParseGlob("templates/*html"))


// Register handler
func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "register.html", nil)
		return
	} else if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Println(email)
	fmt.Println(username)
	fmt.Println(password)
	fmt.Println()

	// Check if user already exists
	_, err := dataBase.GetUserByEmail(email)
	if err == nil {
		http.Error(w, "Email already taken", http.StatusConflict)
		return
	}

	// Create user in the database
	err = dataBase.InsertUser(email, username, password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

