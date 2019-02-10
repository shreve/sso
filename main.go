package main

import (
	"os"
	"log"
	"fmt"
	"time"
	"net/http"

	"github.com/gorilla/mux"
)

type StrMap map[string]string;
type StrListMap map[string][]string;

func getEnv(key, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok { val = def }
	return val
}

var domain = getEnv("AUTH_DOMAIN", "localhost")
var secure = (getEnv("SECURE_ONLY", "true") == "true")

func writeError(w http.ResponseWriter, err error) bool {
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(Render(nil, err))
		return true
	}
	return false
}

func writeCookie(w http.ResponseWriter, user *User) {
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name: "uid",
		Value: user.Uid,
		Domain: domain,
		Expires: expire,
		Secure: secure,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

var root = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write(Render(StrMap{"status": "ok"}, nil))
})

var status = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "uid" {
			user, err := findUserByUid(cookie.Value)
			if writeError(w, err) { return }
			w.Write(Render(genToken(&user), nil))
			return
		}
	}
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

func main() {
	initDB()
	r := mux.NewRouter()

	r.Handle("/", root)
	r.Handle("/status", status)
	r.Handle("/register", register)
	r.Handle("/login", login)

	req := Middleware().Then(r)
	log.Println("================================")
	log.Println("Starting up SSOperhero")
	log.Println("  domain: \t" + domain)
	log.Println("  secure: \t" + fmt.Sprintf("%t", secure))
	log.Println("  db path: \t" + db_path)
	log.Println("  jwt sig: \t" + string(signingKey))
	log.Println("  bcrypt: \t" + fmt.Sprintf("%d", cost))
	log.Println("Serving on http://localhost:9999")
	log.Println("================================")
	log.Println("")
	http.ListenAndServe(":9999", req)
}
