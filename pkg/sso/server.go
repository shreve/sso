package sso

import (
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

type server struct {
	db     *db
	router *mux.Router
	Config *Config
}

func NewServer() *server {
	s := server{}
	s.Config = loadConfig()
	return &s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) ListenAndServe() {
	s.Config.Report()
	s.setupDB()
	s.routes()
	http.ListenAndServe(s.Config.Port, s)
}

var indexView = template.Must(template.ParseFiles("./web/tmpl/index.html"))

func (s *server) signedInUser(r *http.Request) (User, error) {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "uid" {
			return s.db.findUserByUid(cookie.Value)
		}
	}
	return User{}, NotSignedInError
}

func (s *server) clientAllowed(client string) bool {
	for _, allowedClient := range s.Config.Clients {
		if client == allowedClient {
			return true
		}
	}
	return false
}
