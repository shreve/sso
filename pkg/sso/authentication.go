package sso

import (
	"encoding/json"
	"io"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func parseCreds(data io.ReadCloser) (Credentials, error) {
	creds := Credentials{}
	json.NewDecoder(data).Decode(&creds)
	return creds, creds.Validate()
}

func (s *server) hash(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), s.Config.HashCost)
	return string(h), err
}

func compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (s *server) genToken(user *User) Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"source":   s.Config.Domain,
		"username": user.Username,
		"uid":      user.Uid,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokString, _ := token.SignedString(s.Config.TokenSigningKey)
	return Token{tokString, user}
}

func (s *server) loginUser(creds Credentials) (token Token, err error) {
	user, err := s.db.findUserByUsername(creds.Username)
	if err != nil {
		return
	}
	err = compare(user.Password, creds.Password)
	if err != nil {
		return
	}
	return s.genToken(&user), nil
}

func (s *server) registerUser(creds Credentials) (token Token, err error) {
	_, err = s.db.findUserByUsername(creds.Username)
	if err == nil {
		err = AlreadyRegisteredError
		return
	}
	hashed, err := s.hash(creds.Password)
	if err != nil {
		return
	}
	user, err := s.db.createUser(creds.Username, hashed)
	if err != nil {
		return
	}
	return s.genToken(&user), nil
}
