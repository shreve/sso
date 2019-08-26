package main

import (
	"errors"
	"strings"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	"github.com/rs/xid"
)

var db *sql.DB
var db_url = getEnv("DATABASE_URL", "./auth.db")
var db_driver = getEnv("DATABASE_DRIVER", "sqlite3")

type User struct {
	Uid string
	Username string
	Password string `json:"-"`
}

var UserNotFound = errors.New("That user does not exist.")
var UserNotCreated = errors.New("There was a problem creating that user.")

func initDB() {
	var err error
	db, err = sql.Open(db_driver, db_url)
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
		errlog.Println(err)
		return user, UserNotFound
	}
	return user, nil
}

func findUserByUsername(username string) (User, error) {
	result := db.QueryRow(
		"select uid, username, password from users where lower(username) = $1",
		strings.ToLower(username))
	var user User
	err := result.Scan(&user.Uid, &user.Username, &user.Password)
	if err != nil {
		errlog.Println(err)
		return user, UserNotFound
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
		return User{}, UserNotCreated
	}
	return findUserByUid(id.String())
}
