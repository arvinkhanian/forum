package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Post struct {
	ID       int
	Username string
	UserID   int
	Title    string
	Content  string
	Category string
	Likes    int
	Dislikes int
	Comments []comment
}

type comment struct {
	ID       int
	Username string
	UserID   int
	Content  string
	Likes    int
	Dislikes int
}

type postPageData struct {
	Post      Post
	SessionID int
}

type homePageData struct {
	ID       int
	Title    string
	Likes    int
	Dislikes int
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		allPosts, err := fetchPosts()
		if errorCheckHandlers(w, "Failed to load posts", err, http.StatusInternalServerError) {
			return
		}

		// Parse the template
		tmpl, err := template.ParseFiles("./html/home.html")
		if errorCheckHandlers(w, "Failed to parse the template", err, http.StatusInternalServerError) {
			return
		}

		// Execute the template with the posts data
		err = tmpl.Execute(w, allPosts)
		if errorCheckHandlers(w, "Failed to execute the template", err, http.StatusInternalServerError) {
			return
		}
	}
}

// registerHandler handles user registration
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if checkEmailExists(w, email) {
			return
		}

		hashed, err := hashPassword(password)
		if errorCheckHandlers(w, "Password hashing failed", err, http.StatusInternalServerError) {
			return
		}

		if err := saveUser(username, email, hashed); errorCheckHandlers(w, "User registration failed", err, http.StatusInternalServerError) {
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		http.ServeFile(w, r, "./html/register.html")
	}
}

// loginHandler handles user login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		userID, err := authenticateUser(email, password)
		if errorCheckHandlers(w, "Invalid credentials", err, http.StatusUnauthorized) {
			return
		}

		if err := createSession(w, userID); errorCheckHandlers(w, "Session creation failed", err, http.StatusInternalServerError) {
			return
		}

		http.Redirect(w, r, "/home", http.StatusFound)
	} else {
		http.ServeFile(w, r, "./html/login.html")
	}
}

// postHandler handles creating a new post
func createPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userID, err := getUserIDFromSession(r)
		if errorCheckHandlers(w, "Invalid session", err, http.StatusUnauthorized) {
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		category := r.FormValue("category")

		// Insert the post into the database
		if err := createPost(userID, title, content, category); errorCheckHandlers(w, "Post creation failed", err, http.StatusInternalServerError) {
			return
		}

		id, err := getPostId()
		if errorCheckHandlers(w, "Database issue", err, http.StatusInternalServerError) {
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/posts?id=%d", id), http.StatusFound)
	} else {
		fmt.Print("test")
		http.ServeFile(w, r, "./html/createPost.html")
	}
}

// postsHandler displays a single post
func viewPostHandler(w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the URL query parameter
	postID := r.URL.Query().Get("id")
	if postID == "" {
		http.Error(w, "Post ID is missing", http.StatusBadRequest)
		return
	}

	post, err := fetchPost(postID)
	if errorCheckHandlers(w, "Failed to load the post", err, http.StatusInternalServerError) {
		return
	}

	post.Likes, post.Dislikes, err = fetchReactionsNumber(post.ID, false)
	if errorCheckHandlers(w, "Failed to load the reactions number", err, http.StatusInternalServerError) {
		return
	}

	// Fetch comments for this post
	post.Comments, err = fetchCommentsForPost(post.ID)
	if errorCheckHandlers(w, "Failed to load the comments", err, http.StatusInternalServerError) {
		return
	}

	// Pass UserID to the template if logged in
	sessionID, err := getUserIDFromSession(r)
	if err != nil {
		sessionID = 0 // If there's an error, set sessionID to 0
	}

	postPageData := postPageData{
		Post:      post,
		SessionID: sessionID, // Add the user ID
	}

	// Parse the template
	tmpl, err := template.ParseFiles("./html/post.html")
	if errorCheckHandlers(w, "Failed to parse the template", err, http.StatusInternalServerError) {
		return
	}

	// Execute the template, passing in the post data
	err = tmpl.Execute(w, postPageData)
	if errorCheckHandlers(w, "Failed to render the template", err, http.StatusInternalServerError) {
		return
	}
}

// commentHandler handles adding a comment to a post
func commentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	sessionID, err := getUserIDFromSession(r)
	if errorCheckHandlers(w, "Invalid session", err, http.StatusUnauthorized) {
		return
	}

	postID := r.FormValue("post_id")
	if postID == "" {
		http.Error(w, "Post ID is missing", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Content is missing", http.StatusBadRequest)
		return
	}

	if err := addComment(sessionID, postID, content); errorCheckHandlers(w, "Failed to add the comment", err, http.StatusInternalServerError) {
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post?id=%s", postID), http.StatusFound)
}

// likeHandler handles liking or disliking a post or comment
func likeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	sessionID, err := getUserIDFromSession(r)
	if errorCheckHandlers(w, "Invalid session", err, http.StatusUnauthorized) {
		return
	}

	itemID := r.FormValue("item_id")
	if itemID == "" {
		http.Error(w, "Item ID is missing", http.StatusBadRequest)
		return
	}

	isComment := r.FormValue("is_comment") == "true"
	reactionType := r.FormValue("type") // "like" or "dislike"

	if err := likeItem(sessionID, itemID, isComment, reactionType); errorCheckHandlers(w, "Failed to like the item", err, http.StatusInternalServerError) {
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

// logoutHandler handles user logout
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if errorCheckHandlers(w, "No active session", err, http.StatusUnauthorized) {
		return
	}

	if err := deleteSession(cookie.Value); errorCheckHandlers(w, "Logout failed", err, http.StatusInternalServerError) {
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "session_id", Value: "", MaxAge: -1})
	fmt.Fprintf(w, "Logout successful")
	http.Redirect(w, r, "/home", http.StatusFound)
}

// filterHandler handles filtering posts by category
func filterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Get the category from the query parameter
		category := r.URL.Query().Get("category")
		if category == "" {
			http.Error(w, "Category parameter is missing", http.StatusBadRequest)
			return
		}

		// Query the posts for the specified category
		rows, err := db.Query("SELECT id, user_id, title, content, category FROM posts WHERE category = ?", category)
		if err != nil {
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Define a slice to hold the posts
		var posts []struct {
			ID       int
			UserID   int
			Title    string
			Content  string
			Category string
		}

		// Scan the rows into the posts slice
		for rows.Next() {
			var post struct {
				ID       int
				UserID   int
				Title    string
				Content  string
				Category string
			}
			err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category)
			if err != nil {
				http.Error(w, "Failed to scan post", http.StatusInternalServerError)
				return
			}
			posts = append(posts, post)
		}

		// Prepare the data to pass to the template
		data := struct {
			Category string
			Posts    []struct {
				ID       int
				UserID   int
				Title    string
				Content  string
				Category string
			}
		}{
			Category: category,
			Posts:    posts,
		}

		// Parse and execute the template
		tmpl, err := template.ParseFiles("./html/category.html")
		if err != nil {
			http.Error(w, "Failed to parse template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to execute template", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
