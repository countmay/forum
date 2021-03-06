package routes

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"../models"
	"../session"

	"golang.org/x/crypto/bcrypt"
)

// Handler ...
type Handler struct {
	Tmpl            *template.Template
	InMemorySession *session.Session
	Db              *models.Database
}

var posts []models.Post
var users []models.User
var comments []models.Comment
var reactions []models.Reactions

// var s session.SessionData

// CokkieName ...
const (
	CookieName = "sessionID"
)

// IndexHandler ...
func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		fmt.Fprintf(w, "Error 404")
		return
	}
	var authed bool

	cookie, err := r.Cookie("sessionID")
	if err != nil {
		authed = false
	} else {
		authed = h.InMemorySession.Authed(cookie.Value)
	}

	switch r.Method {
	case "GET":
		// fmt.Println(posts)
		PrintPosts()
		model := models.PostListModel{}
		model.IsAuthorized = authed
		model.Posts = posts
		if err := h.Tmpl.ExecuteTemplate(w, "index.html", model); err != nil {
			log.Println(err)
		}
		posts = nil
	}
}

// WriteHandler ...
func (h *Handler) WriteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		model := models.PostListModel{}
		model.IsAuthorized = true
		model.Posts = posts
		h.Tmpl.ExecuteTemplate(w, "write.html", model)
	}

}

// LoginHandler ...
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.Tmpl.ExecuteTemplate(w, "login.html", r)
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")
		fmt.Println(username, password)
		ok, key := h.InMemorySession.CheckUsersSession(username)
		if ok {
			h.InMemorySession.Delete(key)

		}
		if !CheckRepeat(username, "") {
			model := models.PostListModel{}
			model.Exist = false
			h.Tmpl.ExecuteTemplate(w, "login.html", model)
		} else if !CheckPassword(password, username) {
			model := models.PostListModel{}
			model.Exist = false
			h.Tmpl.ExecuteTemplate(w, "login.html", model)

		} else {
			var sessionID string
			sessionID, _ = h.InMemorySession.Init(username)
			// fmt.Println(sessionID)

			cookie := &http.Cookie{
				Name:     CookieName,
				Value:    sessionID,
				Expires:  time.Now().Add(15 * time.Minute),
				HttpOnly: true,
			}

			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusFound)

		}
	}
}

// SigninHandler ...
func (h *Handler) SigninHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.Tmpl.ExecuteTemplate(w, "signin.html", r)
	case "POST":
		userFirstName := r.FormValue("userFirstName")
		userLastName := r.FormValue("userLastName")
		userEmail := r.FormValue("userEmail")
		username := r.FormValue("username")
		password := r.FormValue("password")

		if EmptyMessage(userFirstName) || EmptyMessage(userLastName) || EmptyMessage(userEmail) || EmptyMessage(username) || EmptyMessage(password) {
			model := models.PostListModel{}
			model.EmptyMsg = true
			h.Tmpl.ExecuteTemplate(w, "signin.html", model)
		} else if len(users) > 0 && CheckRepeat(username, userEmail) {
			model := models.PostListModel{}
			model.Exist = true
			h.Tmpl.ExecuteTemplate(w, "signin.html", model)
		} else {
			crypt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				http.Redirect(w, r, "/", 503)
				return
			}

			var sessionID string
			sessionID, _ = h.InMemorySession.Init(username)

			cookie := &http.Cookie{
				Name:     CookieName,
				Value:    sessionID,
				Expires:  time.Now().Add(15 * time.Minute),
				HttpOnly: true,
			}
			u := h.Db.NewUser(userFirstName, userLastName, userEmail, username, crypt, cookie)

			// s.Username = username
			users = append(users, *u)

			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusFound)
		}

	}
}

// LogoutHandler ...
func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		PrintPosts()
		cookie, _ := r.Cookie("sessionID")
		data := h.InMemorySession.Data(cookie.Value)

		ok, key := h.InMemorySession.CheckUsersSession(data.Username)
		if ok {
			h.InMemorySession.Delete(key)
		}
		model := models.PostListModel{}
		// m.IsAuthorized = false
		model.IsAuthorized = false
		model.Posts = posts
		h.Tmpl.ExecuteTemplate(w, "index.html", model)
		posts = nil
	}
}

// SavePostHandler ...
func (h *Handler) SavePostHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("sessionID")
	switch r.Method {
	case "POST":
		title := r.FormValue("titleText")
		categorieLove := r.FormValue("categorieLove")
		categorieFashion := r.FormValue("categorieFashion")
		categorieBeauty := r.FormValue("categorieBeauty")
		categorieHealth := r.FormValue("categorieHealth")
		categoriePopular := r.FormValue("categoriePopular")

		categorie := []string{}
		if categorieLove == "on" {
			categorie = append(categorie, "Любовь")
		}
		if categorieFashion == "on" {
			categorie = append(categorie, "Мода")
		}
		if categorieBeauty == "on" {
			categorie = append(categorie, "Красота")
		}
		if categorieHealth == "on" {
			categorie = append(categorie, "Здоровье")
		}
		if categoriePopular == "on" {
			categorie = append(categorie, "Популярное")
		}
		body := r.FormValue("postText")

		if EmptyMessage(title) || EmptyMessage(body) {
			model := models.PostListModel{}
			model.IsAuthorized = true
			model.Posts = posts
			model.EmptyMsg = true
			h.Tmpl.ExecuteTemplate(w, "write.html", model)
		} else {
			// fmt.Println(title, body, s.Username, categorie)
			data := h.InMemorySession.Data(cookie.Value)
			h.Db.NewPost(title, body, data.Username, categorie, r)

			// posts = append(posts, *post)
			http.Redirect(w, r, "/", 302)
		}

	}
}

// CommentHandler ...
func (h *Handler) CommentHandler(w http.ResponseWriter, r *http.Request) {
	// db, err := sql.Open("sqlite3", "history.db")
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()
	cookie, _ := r.Cookie("sessionID")
	data := h.InMemorySession.Data(cookie.Value)
	switch r.Method {
	case "POST":
		commentText := r.FormValue("comment")
		postUUID := r.FormValue("postUUID")
		// c, _ := r.Cookie("sessionID")
		h.Db.NewComment(postUUID, commentText, data.Username, r)
		// comments = append(comments, *comment)
		http.Redirect(w, r, "/", 302)
	}

}

//ReactionHandler ..
func (h *Handler) ReactionHandler(w http.ResponseWriter, r *http.Request) {

	cook, err := (*http.Request).Cookie(r, "sessionID")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	data := h.InMemorySession.Data(cook.Value)
	fmt.Println(cook)
	switch r.Method {
	case "POST":

		reaction := r.FormValue("reaction")
		postUUID := r.FormValue("postUUID")
		fmt.Println(postUUID)
		fmt.Println(reaction)
		fmt.Println("hello")
		ok, ID := h.Db.CheckReaction(postUUID, data.Username)
		if ok {
			h.Db.DeleteReaction(ID)
		}
		h.Db.NewReaction(data.Username, postUUID, reaction, r)
		http.Redirect(w, r, "/", 302)
	}
}

// PrintPosts ...
func PrintPosts() {
	db, err := sql.Open("sqlite3", "history.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY AUTOINCREMENT,
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

	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		p := models.Post{}
		err := rows.Scan(&p.ID, &p.PostUuid, &p.Username, &p.Title, &p.Categories, &p.Body)
		if err != nil {
			panic(err)
		}

		ReactionArr()
		p.Like, p.Dislike = ReactionCount(p.PostUuid)
		reactions = nil
		p.Comments = PostComments(p.PostUuid)

		posts = append(posts, p)

	}

	reverseArray()
}

// ReactionArr ..
func ReactionArr() {
	db, err := sql.Open("sqlite3", "history.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	statement, err := db.Prepare(`CREATE TABLE IF NOT EXISTS reactions (id INTEGER PRIMARY KEY AUTOINCREMENT,
		username  varchar(128),
		content_uuid varchar(128),
like INteger DEFAULT 0,
dislike integer DEFAULT 0,
		   FOREIGN key (user_uuid) REFERENCES users(uuid),
		   FOREIGN key (content_uuid) REFERENCES posts(post_uuid),
		   FOREIGN KEY (content_uuid) REFERENCES comments(comment_uuid));`)

	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	rows, err := db.Query("SELECT * FROM reactions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		p := models.Reactions{}
		err := rows.Scan(&p.ID, &p.Username, &p.ContentUUID, &p.Like, &p.Dislike)
		if err != nil {
			log.Fatal(err)
		}

		reactions = append(reactions, p)
	}

}

// ReactionCount ..
func ReactionCount(postUUID string) (int, int) {
	likes := 0
	dislikes := 0
	for i := 0; i < len(reactions); i++ {

		if reactions[i].ContentUUID == postUUID {
			if reactions[i].Like == 1 {
				likes++
			}
			if reactions[i].Dislike == 1 {
				dislikes++
			}

		}

	}

	return likes, dislikes
}

// PostComments ..
func PostComments(postUUID string) []models.Comment {
	db, err := sql.Open("sqlite3", "history.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var comments []models.Comment
	rows, err := db.Query("SELECT * FROM comments")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		c := models.Comment{}
		err := rows.Scan(&c.ID, &c.PostUuid, &c.CommentUUID, &c.Author, &c.Body, &c.Like, &c.Dislike)
		if err != nil {
			panic(err)
		}
		ReactionArr()
		c.Like, c.Dislike = ReactionCount(c.CommentUUID)
		reactions = nil
		if c.PostUuid == postUUID {
			comments = append(comments, c)
		}

	}
	return comments
}

func reverseArray() {
	for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
		posts[i], posts[j] = posts[j], posts[i]
	}
}
