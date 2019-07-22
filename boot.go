package gdsm

import (
	"os"
)

// BuildDaemon builds enough of the GDSM Operator so that you can query data internally.
func BuildDaemon() *Operator {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8081"
	}

	operator := New()
	operator.NetAddr = "0.0.0.0:" + port

	return operator
}

// BootDaemon boots both the server and the client.
// Meant to be called `go gdsm.BootDaemon(daemon)` for non blocking and  `gdsm.BootDaemon(daemon)` for blocking.
// If the MANAGER env is set to "true" then the node will not boot a client, it will just be a server for all the workers to attach to.
// If the MANAGER env is not set, the node will boot up a server as well as a client that attaches to the UPLINK server.
func BootDaemon(operator *Operator) {
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

	block()
}

// blocks forever until panic
func block() {
	select {}
}
