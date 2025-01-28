package handlers

import (
"net/http"
//"fmt"
"forum/dataBase"
"forum/structures"
)

func Index(w http.ResponseWriter, r *http.Request) {
    // Render index page with available posts
    rows, err := dataBase.DB.Query("SELECT id, title FROM posts")
    if err != nil {
        http.Error(w, "Error fetching posts", http.StatusInternalServerError)
        return
    }
    defer rows.Close()
	var posts []structures.Post
    for rows.Next() {
        var post structures.Post
        if err := rows.Scan(&post.ID, &post.Title, &post.Content); err != nil {
            http.Error(w, "Error scanning posts", http.StatusInternalServerError)
            return
        }
        posts = append(posts, post)
    }

    data := struct {
        Posts []structures.Post
    }{
        Posts: posts,
    }

    err = templates.ExecuteTemplate(w, "index.html", data)
    if err != nil {
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
    }
}
