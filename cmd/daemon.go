package main

import (
	"os"
	"time"

	"github.com/selfup/jeanome"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8081"
	}

	go func() {
		operator := jeanome.New()
		operator.NetAddr = "127.0.0.1:" + port
		operator.Boot()
	}()

	go func() {
		caller := new(jeanome.Caller)
		caller.NetAddr = "127.0.0.1:" + port
		caller.Dial()
	}()

	recurse()
}

func recurse() {
	time.Sleep(10 * time.Second)
	recurse()
}
