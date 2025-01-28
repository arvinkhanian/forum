package handlers

import (
	"fmt"
	"forum/dataBase"
	"net/http"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        postID := r.FormValue("post_id")
        userID := 1 // Assuming the user is logged in with ID = 1 (You should handle this based on the session)
        content := r.FormValue("content")

        query := "INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)"
        _, err := dataBase.DB.Exec(query, postID, userID, content)
        if err != nil {
            http.Error(w, "Error creating comment", http.StatusInternalServerError)
            return
        }

        fmt.Fprintln(w, "Comment created successfully")
        return
    }

    // Render comment creation form (just an example)
    fmt.Fprintln(w, `<form action="/createComment" method="POST">
                        <textarea name="content" required></textarea><br><br>
                        <input type="hidden" name="post_id" value="1">
                        <button type="submit">Submit Comment</button>
                    </form>`)
}