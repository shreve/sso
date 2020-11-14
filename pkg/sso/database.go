package sso

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/xid"
)

type db struct {
	conn *sql.DB
}

func (s *server) setupDB() {
	conn, err := sql.Open("sqlite3", s.Config.DbUrl)
	if err != nil {
		panic(err)
	}
	conn.Exec(
		"create table if not exists users (" +
			"uid text primary key, " +
			"username text not null, " +
			"password text not null, " +
			"created_at date not null default current_timestamp, " +
			"constraint username_unique unique (username)" +
			")")
	s.db = &db{conn}
}

func (d *db) findUserByUid(uid string) (user User, err error) {
	result := d.conn.QueryRow("select uid, username, password from users where uid = $1", uid)
	err = result.Scan(&user.Uid, &user.Username, &user.Password)
	if err != nil {
		log.Println(err)
		err = UserNotFoundError
	}
	return
}

func (d *db) findUserByUsername(username string) (user User, err error) {
	result := d.conn.QueryRow(
		"select uid, username, password from users where lower(username) = $1",
		strings.ToLower(username))
	err = result.Scan(&user.Uid, &user.Username, &user.Password)
	if err != nil {
		log.Println(err)
		err = UserNotFoundError
	}
	return
}

func (d *db) createUser(username, hashed_password string) (user User, err error) {
	id := xid.New()
	_, err = d.conn.Exec(
		"insert into users (uid, username, password) values ($1, $2, $3)",
		id.String(),
		username,
		hashed_password)
	if err != nil {
		log.Println(err)
		err = UserNotCreatedError
		return
	}
	return d.findUserByUid(id.String())
}
