package sso

import (
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"io"
	"time"
)

func parseCreds(data io.ReadCloser) (Credentials, error) {
	creds := Credentials{}
	json.NewDecoder(data).Decode(&creds)
	return creds, creds.Validate()
}

func hash(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), config.HashCost)
	return string(h), err
}

func compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func genToken(user *User) Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"source":   config.Domain,
		"username": user.Username,
		"uid":      user.Uid,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokString, _ := token.SignedString(config.TokenSigningKey)
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
