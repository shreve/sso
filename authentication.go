package main

import (
	"io"
	"time"
	"errors"
	"strconv"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	jwt "github.com/dgrijalva/jwt-go"
)

var signingKey = []byte(getEnv("JWT_SIGNING_KEY", "jwt signature"))
var cost, err = strconv.Atoi(getEnv("BCRYPT_COST", "10"))

type Credentials struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
}

type Token struct {
	Token string `json:"token"`
	User *User `json:"-"`
}

func (c *Credentials) Validate() error {
	if c.Username == "" || c.Password == "" {
		return EmptyCredsError
	}
	return nil
}

var EmptyCredsError = errors.New("A username and password must be supplied")
var AlreadyRegisteredError = errors.New("That username has already been registered")
var IncorrectPasswordError = errors.New("That password is incorrect")

func parseCreds(data io.ReadCloser) (Credentials, error) {
	creds := Credentials{}
	json.NewDecoder(data).Decode(&creds)
	return creds, creds.Validate()
}

func hash(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(h), err
}

func compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func genToken(user *User) Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"source": domain,
		"username": user.Username,
		"uid": user.Uid,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokString, _ := token.SignedString(signingKey)
	return Token{tokString, user}
}

func loginUser(creds Credentials) (Token, error) {
	user, err := findUserByUsername(creds.Username)
	if err != nil {
		return Token{}, err
	}
	err = compare(user.Password, creds.Password)
	if err != nil {
		return Token{}, IncorrectPasswordError
	}
	return genToken(&user), nil
}

func registerUser(creds Credentials) (Token, error) {
	_, err := findUserByUsername(creds.Username)
	if err == nil {
		return Token{}, AlreadyRegisteredError
	}
	hashed, err := hash(creds.Password)
	if err != nil {
		return Token{}, err
	}
	user, err := createUser(creds.Username, hashed)
	if err != nil {
		return Token{}, err
	}
	return genToken(&user), nil
}
