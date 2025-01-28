package main

import (
	"forum/dataBase"
	//"fmt"
	"forum/handlers"
	"log"
	"net/http"
	//"time"
	//"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	//"golang.org/x/crypto/bcrypt"
	//"database/sql"
	//"html/template"
	//"forum/structures"
)






func main() {


	db, err := dataBase.OpenDatabase()
	if err != nil {
		log.Fatalf("Could not open the database: %v", err)
	}
	defer db.Close()
	dataBase.CreateTables()

	http.Handle("/images/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/home", handlers.Home)
	http.HandleFunc("/createPost", handlers.CreatePost)
    http.HandleFunc("/createComment", handlers.CreateComment)
    //http.HandleFunc("/likePost", handlers.LikePost)
    //http.HandleFunc("/likeComment", handlers.LikeComment)
    http.HandleFunc("/posts", handlers.Posts)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
