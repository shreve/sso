package sso

import (
	"fmt"
	"log"
	"crypto/md5"
)

type Config struct {
	Domain string
	Clients []string
	SecureCookies bool
	HashCost int
	TokenSigningKey []byte
	DbUrl string
}

var config *Config

func (c *Config) Report() {
	log.Println("================================")
	log.Println("Starting up SSOperhero")
	log.Println("  host domain: \t" + c.Domain)
	log.Println("  secure cookies: \t" + fmt.Sprintf("%t", c.SecureCookies))
	log.Println("  database: \t" + c.DbUrl)
	log.Println("  bcrypt rounds: \t" + fmt.Sprintf("%d", c.HashCost))
	log.Println("  signature hash: \t" + fmt.Sprintf("%x", md5.Sum(c.TokenSigningKey)))
	log.Println("  clients: \t" + fmt.Sprintf("%v", c.Clients))
	log.Println("================================")
}
