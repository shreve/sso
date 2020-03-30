package main

import (
	"os"
	"log"
	"strconv"
	"strings"
	"net/http"

	"github.com/shreve/sso"
)

func getEnv(key, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok { val = def }
	return val
}

func main() {
	cost, err := strconv.Atoi(getEnv("BCRYPT_COST", "10"))
	if err != nil {
		log.Print("Unable to parse BCRYPT_COST parameter. Using default 10.")
		cost = 10
	}

	key, ok := os.LookupEnv("JWT_SIGNING_KEY")
	if !ok {
		log.Fatal("Can't run server without a signing key JWT_SIGNING_KEY.")
	}

	config := sso.Config{
		Domain: getEnv("AUTH_DOMAIN", "localhost"),
		Clients: strings.Split(getEnv("CLIENT_DOMAINS", ""), ","),
		SecureCookies: "true" == getEnv("SECURE_ONLY", "true"),
		HashCost: cost,
		TokenSigningKey: []byte(key),
		DbUrl: getEnv("DATABASE_URL", "./auth.db"),
	}

	config.Report()
	mux := sso.NewServer(&config)

	port := getEnv("PORT", "9999")
	log.Println("Starting server on :" + port)
	http.ListenAndServe(":" + port, mux)
}
