package sso

type User struct {
	Uid      string
	Username string
	Password string `json:"-"`
}
