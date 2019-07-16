package gdsm

import (
	"os"
	"time"
)

// BootMattDaemon boots both the server and the client
// Meant to be called: go BootMattDaemon()
func BootMattDaemon() {
	port := os.Getenv("PORT")
	netAddr := os.Getenv("UPLINK")
	manager := os.Getenv("MANAGER")

	if port == "" {
		port = "8081"
	}

	operator := New()
	operator.NetAddr = "0.0.0.0:" + port

	go func() {
		operator.Boot()
	}()

	// if this node is a manager then do not have the client connect to self
	if manager == "" {
		go func() {
			caller := new(Caller)

			if netAddr != "" {
				caller.NetAddr = netAddr
				caller.Server = operator.NetAddr
			} else {
				caller.NetAddr = "0.0.0.0:" + port
				caller.Server = operator.NetAddr
			}

			caller.Dial()
		}()
	}

	recurse()
}

func recurse() {
	time.Sleep(1 * time.Second)
	recurse()
}
