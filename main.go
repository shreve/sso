package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type StrMap map[string]string;
type StrListMap map[string][]string;

func writeError(w http.ResponseWriter, err error) bool {
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(Render(nil, err))
		return true
	}
	return false
}

var root = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write(Render(StrMap{"status": "ok"}, nil))
})

var register = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	creds, err := parseCreds(r.Body)
	if writeError(w, err) { return }
	user, err := registerUser(creds)
	if writeError(w, err) { return }
	w.Write(Render(user, nil))
})

var login = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	creds, err := parseCreds(r.Body)
	if writeError(w, err) { return }
	user, err := loginUser(creds)
	if writeError(w, err) { return }
	w.Write(Render(user, nil))
})

func main() {
	initDB()
	r := mux.NewRouter()

	r.Handle("/", root)
	r.Handle("/register", register)
	r.Handle("/login", login)

	req := Middleware().Then(r)
	log.Println("Serving on http://localhost:9999")
	http.ListenAndServe(":9999", req)
}
