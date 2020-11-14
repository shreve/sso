package sso

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

func (s *server) routes() {
	s.router = mux.NewRouter()
	s.router.Use(s.logging)
	s.router.Use(s.cors)

	s.router.PathPrefix("/js").Handler(http.FileServer(http.Dir("./web/")))

	s.router.HandleFunc("/", s.root).Methods("GET")

	sub := s.router.PathPrefix("/").Subrouter()
	sub.Use(s.contentType("application/json"))

	sub.HandleFunc("/status", s.status).Methods("GET")
	sub.HandleFunc("/register", s.register).Methods("POST")
	sub.HandleFunc("/login", s.login).Methods("POST")
	sub.HandleFunc("/logout", s.logout).Methods("GET", "POST")
}

func (s *server) root(w http.ResponseWriter, r *http.Request) {
	// Check if this is from an allowed domain
	client, err := url.Parse(r.Header.Get("Referer"))
	if !s.clientAllowed(client.String()) {
		return
	}

	// s.render the domain into the script-loader index file
	if !s.writeError(w, err) {
		w.Header().Set("Content-Type", "text/html")
		err = indexView.Execute(w, struct{ Domain string }{client.String()})
		s.writeError(w, err)
	}
}

func (s *server) status(w http.ResponseWriter, r *http.Request) {
	// If there's a signed in error, fetch a new token for them
	user, err := s.signedInUser(r)
	if s.writeError(w, err) {
		return
	}
	w.Write(s.render(s.genToken(&user), nil))
}

func (s *server) register(w http.ResponseWriter, r *http.Request) {
	creds, err := parseCreds(r.Body)
	if s.writeError(w, err) {
		return
	}
	token, err := s.registerUser(creds)
	if s.writeError(w, err) {
		return
	}
	s.writeCookie(w, "uid", token.User.Uid)
	w.Write(s.render(token, nil))
}

func (s *server) login(w http.ResponseWriter, r *http.Request) {
	creds, err := parseCreds(r.Body)
	if s.writeError(w, err) {
		return
	}
	token, err := s.loginUser(creds)
	if s.writeError(w, err) {
		return
	}
	s.writeCookie(w, "uid", token.User.Uid)
	w.Write(s.render(token, nil))
}

func (s *server) logout(w http.ResponseWriter, r *http.Request) {
	s.clearCookie(w, "uid")
}
