package sso

import (
	"net/url"
	"net/http"
	"text/template"
	"github.com/gorilla/mux"
)

func root(w http.ResponseWriter, r *http.Request) {
	// Check if this is from an allowed domain
	client, err := url.Parse(r.Header.Get("Referer"))
	if ! clientAllowed(client.Host) { return }

	// Render the domain into the script-loader index file
	t, err := template.ParseFiles("../web/tmpl/index.html")
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, struct { Domain string }{client.Host})
	writeError(w, err)
}

func status(w http.ResponseWriter, r *http.Request) {
	// If there's a signed in error, fetch a new token for them
	user, err := signedInUser(r)
	if writeError(w, err) { return }
	w.Write(render(genToken(&user), nil))
}

func register(w http.ResponseWriter, r *http.Request) {
	creds, err := parseCreds(r.Body)
	if writeError(w, err) { return }
	token, err := registerUser(creds)
	if writeError(w, err) { return }
	writeCookie(w, "uid", token.User.Uid)
	w.Write(render(token, nil))
}

func login(w http.ResponseWriter, r *http.Request) {
	creds, err := parseCreds(r.Body)
	if writeError(w, err) { return }
	token, err := loginUser(creds)
	if writeError(w, err) { return }
	writeCookie(w, "uid", token.User.Uid)
	w.Write(render(token, nil))
}

func logout(w http.ResponseWriter, r *http.Request) {
	clearCookie(w, "uid")
}

func signedInUser(r *http.Request) (User, error){
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "uid" {
			return findUserByUid(cookie.Value);
		}
	}
	return User{}, NotSignedInError
}

func clientAllowed(client string) bool {
	for _, allowedClient := range config.Clients {
		if client == allowedClient { return true }
	}
	return false
}

func NewServer(c *Config) *mux.Router {
	config = c

	initializeDB(config.DbUrl)

	r := mux.NewRouter()
	r.Use(Logging)
	r.Use(CORS)

	r.PathPrefix("/js").Handler(http.FileServer(http.Dir("../web/")))

	s := r.PathPrefix("/").Subrouter()
	s.Use(ContentType("application/json"))

	s.HandleFunc("/", root)
	s.HandleFunc("/status", status)
	s.HandleFunc("/register", register)
	s.HandleFunc("/login", login)
	s.HandleFunc("/logout", logout)

	return r
}
