package handlers

import (
	"net/http"
	"fmt"
	"forum/dataBase"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Assuming session ID is already handled here
    userID := getUserIDFromSession(r) // You would need to implement session validation
    title := r.FormValue("title")
    content := r.FormValue("content")

    // Insert the post into the database
    _, err := dataBase.DB.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, title, content)
    if err != nil {
        http.Error(w, "Error creating post", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/home", http.StatusSeeOther)
    fmt.Fprintln(w, "Post created successfully")
}

func getUserIDFromSession(r *http.Request) int {
    cookie, err := r.Cookie("session_id")
    if err != nil {
        return 0 // Session doesn't exist
    }

    var userID int
    err = dataBase.DB.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", cookie.Value).Scan(&userID)
    if err != nil {
        return 0 // Invalid session
    }

    return userID
}
