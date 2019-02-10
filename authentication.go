package main

import (
	"io"
	"os"
	"time"
	"errors"
	"strconv"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	jwt "github.com/dgrijalva/jwt-go"
)

var signingKey = []byte(os.Getenv("JWT_SIGNING_KEY"))
var cost, err = strconv.Atoi(os.Getenv("BCRYPT_COST"))

type Credentials struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
}

type Token struct {
	Token string `json:"token"`
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

func genToken(user *User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"uid": user.Uid,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokString, _ := token.SignedString(signingKey)
	return tokString
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
	return Token{genToken(&user)}, nil
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
	return Token{genToken(&user)}, nil
}
