package gdsm

import (
	"net"
	"strings"
)

// Read  handles payloads without params
func (op *Operator) Read(newMessage string, conn net.Conn) {
	switch newMessage {
	case "ping":
		conn.Write([]byte("pong\n"))
	case "clients":
		clients := op.ConnectedClients()
		clientsStr := strings.Join(clients, "|")
		conn.Write([]byte(clientsStr + "\n"))
	case "nodes", "servers":
		servers := op.Nodes()
		serversStr := strings.Join(servers, "|")
		conn.Write([]byte(serversStr + "\n"))
	case "workers":
		workers := op.Workers()
		workersStr := strings.Join(workers, "|")
		conn.Write([]byte(workersStr + "\n"))
	default:
		conn.Write([]byte("unknown command\n"))
	}
}
