package sso

import (
	"log"
	"strings"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	"github.com/rs/xid"
)

var db *sql.DB

func initializeDB(url string) {
	var err error
	db, err = sql.Open("sqlite3", url)
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
		log.Println(err)
		return user, UserNotFoundError
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
		log.Println(err)
		return user, UserNotFoundError
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
		return User{}, UserNotCreatedError
	}
	return findUserByUid(id.String())
}
