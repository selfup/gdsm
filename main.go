package main

import (
	"os"
	"time"

	"github.com/selfup/gdsm/gdsm"
)

func main() {
	port := os.Getenv("PORT")
	netAddr := os.Getenv("UPLINK")

	if port == "" {
		port = "8081"
	}

	operator := gdsm.New()
	operator.NetAddr = "0.0.0.0:" + port

	go func() {
		operator.Boot()
	}()

	go func() {
		caller := new(gdsm.Caller)

		if netAddr != "" {
			caller.NetAddr = netAddr
			caller.Server = operator.NetAddr
		} else {
			caller.NetAddr = "0.0.0.0:" + port
			caller.Server = operator.NetAddr
		}

		caller.Dial()
	}()

	recurse()
}

func recurse() {
	time.Sleep(1 * time.Second)
	recurse()
}
