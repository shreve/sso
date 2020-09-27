package sso

type Token struct {
	Token string `json:"token"`
	User  *User  `json:"-"`
}
