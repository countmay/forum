package routes

import (
	"database/sql"

	"../models"
	"golang.org/x/crypto/bcrypt"
)

// EmptyMessage ...
func EmptyMessage(text string) bool {
	cnt := 0
	for _, s := range text {
		if s == ' ' || s == '\n' || s == '\r' {
			cnt++
		}
	}

	if cnt == len(text) {
		return true
	}
	return false
}

// CheckRepeat ...
func CheckRepeat(name, email string) bool {
	db, err := sql.Open("sqlite3", "history.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		u := models.User{}
		err := rows.Scan(&u.ID, &u.Uuid, &u.UserFirstName, &u.UserLastName, &u.UserEmail, &u.Username, &u.Password, &u.Cookie)
		if err != nil {
			panic(err)
		}
		if u.Username == name || u.UserEmail == email {
			return true
		}
	}
	return false
}

func CheckPassword(password, login string) bool {
	db, err := sql.Open("sqlite3", "history.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		u := models.User{}
		err := rows.Scan(&u.ID, &u.Uuid, &u.UserFirstName, &u.UserLastName, &u.UserEmail, &u.Username, &u.Password, &u.Cookie)
		if err != nil {
			panic(err)
		}
		if u.Username == login {
			return CheckPasswordHash(password, u.Password)

		}
	}
	return false

}

func CheckPasswordHash(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
