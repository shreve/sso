package main

import (
	"os"
	"errors"
	"strings"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/xid"
)

func getEnv(key, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok { val = def }
	return val
}

var db *sql.DB
var db_path = getEnv("DATABASE_PATH", "./auth.db")

type User struct {
	Uid string
	Username string
	Password string
}

var UserNotFound = errors.New("That username has not been claimed.")
var UserNotCreated = errors.New("There was a problem creating that user.")

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", db_path)
	if err != nil {
		panic(err)
	}
	createDB()
}

func createDB() {
	db.Exec(
		"create table if not exists users (" +
			"uid text primary key, " +
			"username text not null, " +
			"password text not null, " +
			"created_at date not null default current_timestamp, " +
			"constraint username_unique unique (username)" +
		")")
}

func findUserByUid(uid string) (User, error) {
	result := db.QueryRow("select uid, username, password from users where uid = $1", uid)
	var user User
	err := result.Scan(&user.Uid, &user.Username, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}

func findUserByUsername(username string) (User, error) {
	result := db.QueryRow(
		"select uid, username, password from users where lower(username) = $2",
		strings.ToLower(username))
	var user User
	err := result.Scan(&user.Uid, &user.Username, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}

func createUser(username, hashed_password string) (User, error) {
	id := xid.New()
	_, err := db.Exec(
		"insert into users (uid, username, password) values ($1, $2, $3)",
		id.String(),
		username,
		hashed_password)
	if err != nil {
		return User{}, err
	}
	return findUserByUid(id.String())
}
