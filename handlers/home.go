package handlers

import (
    "fmt"
    "net/http"
  //  "html/template"
   // "log"
    "forum/structures"
	"forum/dataBase"
)

func Home(w http.ResponseWriter, r *http.Request) {
    posts, err := getPosts() // Fetch posts from the DB
    if err != nil {
        http.Error(w, "Could not retrieve posts", http.StatusInternalServerError)
        return
    }

    // Create a map for template data
    data := map[string]interface{}{
        "Posts": posts,
    }

	fmt.Println("test: submitting a post")
	fmt.Println(data)

    // Render the HTML template
	templates.ExecuteTemplate(w, "home.html", data)
}

// Function to fetch posts from the database
func getPosts() ([]structures.Post, error) {
    rows, err := dataBase.DB.Query("SELECT id, title, content FROM posts ORDER BY created_at DESC")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []structures.Post
    for rows.Next() {
        var post structures.Post
        err := rows.Scan(&post.ID, &post.Title, &post.Content)
        if err != nil {
            return nil, err
        }
        posts = append(posts, post)
    }

    return posts, nil
}

