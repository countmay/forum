package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Conn *sql.DB
}

// Post ...
type Post struct {
	ID       int
	PostUuid string
	Username string
	// Author     string
	Title      string
	Categories string
	Body       string
	Comments   []Comment
	Like       int
	Dislike    int
	// Comments   []*Comment
}

// type Comment struct {
// 	CommentID int
// 	Author    string
// 	Body      string
// 	Like      int
// 	Dislike   int
// }

// User ...
type User struct {
	ID            int
	Uuid          string
	UserFirstName string
	UserLastName  string
	UserEmail     string
	Username      string
	Password      []byte
	Cookie        string
}

// Comment ..
type Comment struct {
	ID          int
	PostUuid    string
	CommentUUID string
	Author      string
	Body        string
	Like        int
	Dislike     int
}
type Reactions struct {
	ID          int
	Username    string
	ContentUUID string
	Like        int
	Dislike     int
}

// NewUser ...
func (db *Database) NewUser(userFirstName, userLastName, userEmail, username string, password []byte, cookie *http.Cookie) *User {

	newID, _ := uuid.NewV4()

	statement, err := db.Conn.Prepare(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT,
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

	statement, err = db.Conn.Prepare("INSERT INTO users(uuid, fname, lname, email, login, password, cookie) VALUES(?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		panic(err)
	}
	statement.Exec(newID.String(), userFirstName, userLastName, userEmail, username, password, cookie.Value)
	rows, err := db.Conn.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.ID, &u.Uuid, &u.UserFirstName, &u.UserLastName, &u.UserEmail, &u.Username, &u.Password, &u.Cookie)
		if err != nil {
			panic(err)
		}
		if u.Username == username {
			return &User{u.ID, u.Uuid, u.UserFirstName, u.UserLastName, u.UserEmail, u.Username, u.Password, u.Cookie}
		}
	}
	return nil
}

// NewPost ...
func (db *Database) NewPost(title, body, author string, categories []string, r *http.Request) {

	newID, _ := uuid.NewV4()
	statement, err := db.Conn.Prepare(`CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_uuid varchar(128),
		username  varchar(128),
		title varchar(255),
		categories varchar(255),
		body TEXT,
		FOREIGN KEY (user_uuid) REFERENCES users(uuid));`)
	if err != nil {
		panic(err)
	}
	statement.Exec()
	statement, err = db.Conn.Prepare("INSERT INTO posts(post_uuid, username, title, categories, body) VALUES(?, ?, ?, ?, ?);")
	if err != nil {
		panic(err)
	}
	categorie := ""
	for index, i := range categories {
		if index != len(categories)-1 {
			categorie += i + ","
		} else {
			categorie += i
		}
	}
	// cook, _ := r.Cookie("sessionID")
	// var detectedUser string
	// query := db.Conn.QueryRow(`select uuid from users where cookie = ?`, cook.Value)
	// query.Scan(&detectedUser)
	// name := ""
	// n := db.Conn.QueryRow(`select login from users where uuid =?`, detectedUser)
	// n.Scan(&name)

	// fmt.Println("name:", name)
	statement.Exec(newID.String(), author, title, categorie, body)

}

// NewComment ..
func (db *Database) NewComment(postUUID, commentText, author string, r *http.Request) {
	newID, _ := uuid.NewV4()

	fmt.Println(author, postUUID)
	statement, err := db.Conn.Prepare("INSERT INTO comments(post_uuid,comment_uuid,username,body,like,dislike) VALUES(?, ?, ?, ?, ?, ?);")
	if err != nil {
		panic(err)
	}

	// cook, _ := r.Cookie("sessionID")
	// var detectedUser string
	// query := db.Conn.QueryRow(`select uuid from users where cookie = ?`, cook.Value)
	// query.Scan(&detectedUser)

	statement.Exec(postUUID, newID, author, commentText, 0, 0)

}

// NewReaction ..
func (db *Database) NewReaction(username, contentUuid, reaction string, r *http.Request) {

	statement, err := db.Conn.Prepare(`CREATE TABLE IF NOT EXISTS reactions (id INTEGER PRIMARY KEY AUTOINCREMENT,
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

	fmt.Println("hey3")
	statement.Exec()

	fmt.Println("hey2")
	statement, err = db.Conn.Prepare("INSERT INTO reactions(username,content_uuid,like,dislike) VALUES(?, ?, ?, ?);")
	if err != nil {
		log.Fatal(err)
	}

	// cook, _ := r.Cookie("sessionID")
	// var detectedUser string
	// query := db.Conn.QueryRow(`select uuid from users where cookie = ?`, cook.Value)
	// query.Scan(&detectedUser)
	arr := strings.Split(reaction, "_")
	res1, err := strconv.Atoi(arr[0])
	if err != nil {
		panic(err)
	}
	res2, err := strconv.Atoi(arr[1])
	if err != nil {
		panic(err)
	}

	statement.Exec(username, contentUuid, res1, res2)

}

func (db *Database) CheckReaction(contentUUID, username string) (bool, int) {
	rows, err := db.Conn.Query("SELECT * FROM reactions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		u := Reactions{}
		err := rows.Scan(&u.ID, &u.Username, &u.ContentUUID, &u.Like, &u.Dislike)
		if err != nil {
			panic(err)
		}
		if (u.Username == username) && (contentUUID == u.ContentUUID) {
			return true, u.ID
		}
	}
	return false, 0
}

func (db *Database) DeleteReaction(ID int) {
	// delete
	statement, err := db.Conn.Prepare("delete from reactions where id=?")

	if err != nil {
		log.Fatal(err)
	}
	statement.Exec(ID)
}

func fn0(err error) {
	if err != nil {
		panic(err)
	}
}
