package sso

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func (s *server) writeError(w http.ResponseWriter, err error) bool {
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(s.render(nil, err))
		return true
	}
	return false
}

func (s *server) writeCookie(w http.ResponseWriter, name, value string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Domain:   s.Config.Domain,
		Expires:  time.Now().AddDate(0, 0, 1),
		Secure:   s.Config.SecureCookies,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func (s *server) clearCookie(w http.ResponseWriter, name string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    "",
		Domain:   s.Config.Domain,
		Expires:  time.Unix(0, 0),
		Secure:   s.Config.SecureCookies,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

func (s *server) errJSON(msg string) []byte {
	return []byte("{\"error\": \"" + msg + "\"}")
}

func (s *server) render(obj interface{}, err error) []byte {
	if err != nil {
		return s.errJSON(err.Error())
	}
	payload, err := json.Marshal(obj)
	if err != nil {
		return s.errJSON(err.Error())
	}
	return payload
}

func boxPrint(lines []string) {
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}

	border := strings.Repeat("=", width+1)
	log.Println(border)
	for _, line := range lines {
		log.Println(line)
	}
	log.Println(border)
}
