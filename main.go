package main

import (
	"os"
	"log"
	"fmt"
	"time"
	"errors"
	"strings"
	"net/url"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func getEnv(key, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok { val = def }
	return val
}

var port = getEnv("PORT", "9999")
var domain = getEnv("AUTH_DOMAIN", "localhost")
var clients = strings.Split(getEnv("CLIENT_DOMAINS", ""), ",")
var secure = (getEnv("SECURE_ONLY", "true") == "true")
var NotSignedIn = errors.New("There is not a signed in user.")
var errlog = log.New(os.Stdout, "error", log.LstdFlags)

type Client struct {
	Domain string
	Allowed bool
}

func writeError(w http.ResponseWriter, err error) bool {
	if err != nil {
		errlog.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(Render(nil, err))
		return true
	}
	return false
}

func writeCookie(w http.ResponseWriter, user *User) {
	cookie := http.Cookie{
		Name: "uid",
		Value: user.Uid,
		Domain: domain,
		Expires: time.Now().AddDate(0, 0, 1),
		Secure: secure,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func clearCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name: "uid",
		Value: "",
		Domain: domain,
		Expires: time.Unix(0, 0),
		Secure: secure,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func signedInUser(r *http.Request) (User, error){
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "uid" {
			return findUserByUid(cookie.Value);
		}
	}
	return User{}, NotSignedIn
}

func getClient(domain string) Client {
	c := Client{"", false}
	u, err := url.Parse(domain)
	if err != nil { return c }
	for _, client := range clients {
		if u.Host == client {
			c.Allowed = true
			c.Domain = (&url.URL{
				Scheme: u.Scheme,
				Host: u.Host,
			}).String()
		}
	}
	return c
}

var root = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	client := getClient(r.Header.Get("Referer"))
	if ! client.Allowed { return }
	t, err := template.ParseFiles("tmpl/index.html")
	// user, err := signedInUser(r)
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, client)
	if writeError(w, err) { return }
})

var status = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user, err := signedInUser(r)
	if writeError(w, err) { return }
	w.Write(Render(genToken(&user), nil))
})

var register = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	creds, err := parseCreds(r.Body)
	if writeError(w, err) { return }
	token, err := registerUser(creds)
	if writeError(w, err) { return }
	writeCookie(w, token.User)
	w.Write(Render(token, nil))
})

var login = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	creds, err := parseCreds(r.Body)
	if writeError(w, err) { return }
	token, err := loginUser(creds)
	if writeError(w, err) { return }
	writeCookie(w, token.User)
	w.Write(Render(token, nil))
})

var logout = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	clearCookie(w)
})

func main() {
	initDB()
	r := mux.NewRouter()
	r.Use(Logging)
	r.Use(CORS)
	r.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./")))

	s := r.PathPrefix("/").Subrouter()
	s.Use(ForceJSON)

	s.Handle("/", root)
	s.Handle("/status", status)
	s.Handle("/register", register)
	s.Handle("/login", login)
	s.Handle("/logout", logout)


	log.Println("================================")
	log.Println("Starting up SSOperhero")
	log.Println("  domain: \t" + domain)
	log.Println("  secure: \t" + fmt.Sprintf("%t", secure))
	log.Println("  db: \t" + db_driver + " " + db_url)
	log.Println("  jwt sig: \t" + string(signingKey))
	log.Println("  bcrypt: \t" + fmt.Sprintf("%d", cost))
	log.Println("Serving on http://localhost:" + port)
	log.Println("================================")
	log.Println("")
	http.ListenAndServe(":" + port, r)
}
