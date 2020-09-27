package sso

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func writeError(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(render(nil, err))
		return true
	}
	return false
}

func writeCookie(w http.ResponseWriter, name, value string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Domain:   config.Domain,
		Expires:  time.Now().AddDate(0, 0, 1),
		Secure:   config.SecureCookies,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func clearCookie(w http.ResponseWriter, name string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    "",
		Domain:   config.Domain,
		Expires:  time.Unix(0, 0),
		Secure:   config.SecureCookies,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func errJSON(msg string) []byte {
	return []byte("{\"error\": \"" + msg + "\"}")
}

func render(obj interface{}, err error) []byte {
	if err != nil {
		return errJSON(err.Error())
	}
	payload, err := json.Marshal(obj)
	if err != nil {
		return errJSON(err.Error())
	}
	return payload
}
