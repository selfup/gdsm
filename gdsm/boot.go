package gdsm

import (
	"os"
	"time"
)

// BuildGDSM builds enough of the GDSM Operator so that you can query data internally
func BuildGDSM() *Operator {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8081"
	}

	operator := New()
	operator.NetAddr = "0.0.0.0:" + port

	return operator
}

// BootMattDaemon boots both the server and the client
// Meant to be called: go BootMattDaemon()
func BootMattDaemon(operator *Operator) {
	netAddr := os.Getenv("UPLINK")
	manager := os.Getenv("MANAGER")
	port := os.Getenv("PORT")

	if port == "" {
		port = "8081"
	}

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
