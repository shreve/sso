package main

import (
	"strings"
	"errors"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type User struct {
	Username string
	Password string
}

var UserNotFound = errors.New("That username has not been claimed.")
var UserNotCreated = errors.New("There was a problem creating that user.")

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "/auth.db")
	if err != nil {
		panic(err)
	}
	createDB()
}

func createDB() {
	db.Exec(
		"create table if not exists users (" +
			"username text primary key, " +
			"password text)")
}

func findUser(username string) (User, error) {
	result := db.QueryRow(
		"select username, password from users where lower(username) = $1",
		strings.ToLower(username))
	var user User
	err := result.Scan(&user.Username, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}

func createUser(username, hashed_password string) (User, error) {
	_, err := db.Exec("insert into users values ($1, $2)", username, hashed_password)
	if err != nil {
		return User{}, err
	}
	return findUser(username)
}
