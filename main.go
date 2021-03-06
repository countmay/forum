package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"./models"
	"./routes"
	"./session"
)

// InMemorySession ...
// var InMemorySession *session.Session

func main() {
	db, err := sql.Open("sqlite3", "history.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid varchar(128),
		  fname varchar(100),
		  lname varchar(50),
		  email varchar(255),
		  login varchar(255),
		  password varchar(255),
		  cookie varchar(255));`)
	if err != nil {
		panic(err)
	}
	statement.Exec()
	statement, err = db.Prepare(`CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_uuid varchar(128),
		username  varchar(128),
		title varchar(255),
		categories varchar(255),
		body TEXT,
		FOREIGN KEY (username) REFERENCES users(login));`)
	if err != nil {
		panic(err)
	}
	statement.Exec()
	statement, err = db.Prepare(`CREATE TABLE IF NOT EXISTS comments (id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_uuid varchar(128),
		comment_uuid varchar(128),
		username varchar(128) NOT NULL,
		body TEXT not NULL,
		like INteger,
		dislike integer,
		FOREIGN key (username) REFERENCES users(login),
		FOREIGN KEY (post_uuid) REFERENCES posts(post_uuid));`)

	if err != nil {
		panic(err)
	}
	statement.Exec()

	statement, err = db.Prepare(`CREATE TABLE IF NOT EXISTS reactions (id INTEGER PRIMARY KEY AUTOINCREMENT,
		username  varchar(128),
		content_uuid varchar(128),
like INteger DEFAULT 0,
dislike integer DEFAULT 0,
		   FOREIGN key (username) REFERENCES users(login),
		   FOREIGN key (content_uuid) REFERENCES posts(post_uuid),
		   FOREIGN KEY (content_uuid) REFERENCES comments(comment_uuid));`)

	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
	database := &models.Database{
		Conn: db,
	}
	handler := routes.Handler{
		Tmpl:            template.Must(template.ParseGlob("./static/templates/*")),
		InMemorySession: session.NewSession(),
		Db:              database,
	}

	router := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	router.HandleFunc("/", handler.IndexHandler)
	router.HandleFunc("/write", handler.AuthMiddleware(handler.WriteHandler))
	router.HandleFunc("/SavePost", handler.AuthMiddleware(handler.SavePostHandler))
	router.HandleFunc("/login", handler.LoginHandler)
	router.HandleFunc("/signin", handler.SigninHandler)
	router.HandleFunc("/logout", handler.LogoutHandler)
	router.HandleFunc("/comment", handler.AuthMiddleware(handler.CommentHandler))
	router.HandleFunc("/reactions", handler.AuthMiddleware(handler.ReactionHandler))

	server := new(http.Server)

	server.Addr = ":8000"
	server.Handler = router
	fmt.Println("Listening on port : 8000")
	server.ListenAndServe()
}
