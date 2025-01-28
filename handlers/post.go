package handlers

import ("net/http"
"forum/structures"
"forum/dataBase"
"fmt"
)

// Filter posts (by categories, likes, etc.)
func Posts(w http.ResponseWriter, r *http.Request) {
    category := r.URL.Query().Get("category")
    liked := r.URL.Query().Get("liked")

    query := "SELECT id, title, content FROM posts WHERE 1=1"
    if category != "" {
        query += " AND id IN (SELECT post_id FROM post_categories WHERE category_id = (SELECT id FROM categories WHERE name = ?))"
    }
    if liked == "true" {
        query += " AND likes > dislikes"
    }

    rows, err := dataBase.DB.Query(query, category)
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

    // Render posts (implement rendering here)
    fmt.Fprintf(w, "Filtered Posts: %+v", data)
}
