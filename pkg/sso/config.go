package sso

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Domain          string
	Clients         []string
	SecureCookies   bool
	HashCost        int
	TokenSigningKey []byte
	DbUrl           string
	Port            string
}

func getEnv(key, def string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		val = def
	}
	return val
}

func loadConfig() *Config {

	cost, err := strconv.Atoi(getEnv("AUTH_BCRYPT_COST", "10"))
	if err != nil {
		log.Println("Unable to parse AUTH_BCRYPT_COST parameter. Using default 10.")
		cost = 10
	}

	key, ok := os.LookupEnv("AUTH_SIGNING_KEY")
	if !ok {
		log.Fatal("Can't run server without a signing key AUTH_SIGNING_KEY.")
	}

	clients := strings.Split(getEnv("AUTH_CLIENT_DOMAINS", ""), ",")
	if len(clients) == 0 {
		log.Fatal("Without clients, the server will not work for anyone.")
	}

	config := Config{
		Port:            getEnv("AUTH_PORT", ":9999"),
		Domain:          getEnv("AUTH_DOMAIN", "localhost"),
		DbUrl:           getEnv("AUTH_DATABASE_URL", "./auth.db"),
		Clients:         clients,
		SecureCookies:   "true" == getEnv("AUTH_SECURE_ONLY", "true"),
		HashCost:        cost,
		TokenSigningKey: []byte(key),
	}

	return &config

}

func (c *Config) Report() {
	boxPrint([]string{
		"Starting with Config",
		"  host domain: \t" + c.Domain,
		"  host port: \t" + c.Port,
		"  secure cookies: \t" + fmt.Sprintf("%t", c.SecureCookies),
		"  database: \t" + c.DbUrl,
		"  bcrypt rounds: \t" + fmt.Sprintf("%d", c.HashCost),
		"  signature hash: \t" + fmt.Sprintf("%x", md5.Sum(c.TokenSigningKey)),
		"  clients: \t\t" + fmt.Sprintf("%v", c.Clients),
	})
}
