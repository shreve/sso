package main

import (
	"github.com/shreve/sso/pkg/sso"
)

func main() {
	sso.NewServer().ListenAndServe()
}
