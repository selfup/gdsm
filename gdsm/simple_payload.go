package gdsm

import (
	"net"
	"strings"
)

// HandleSimplePayload  handles payloads without params
func (op *Operator) HandleSimplePayload(newMessage string, conn net.Conn) {
	switch newMessage {
	case "ping":
		conn.Write([]byte("200\n"))
		break
	case "clients":
		clients := op.ConnectedClients()
		clientsStr := strings.Join(clients, "|")
		conn.Write([]byte(clientsStr + "\n"))
		break
	case "nodes", "servers":
		servers := op.Nodes()
		serversStr := strings.Join(servers, "|")

		conn.Write([]byte(serversStr + "\n"))
		break
	case "workers":
		workers := op.Workers()
		workersStr := strings.Join(workers, "|")

		conn.Write([]byte(workersStr + "\n"))
		break
	default:
		conn.Write([]byte("405\n"))
		break
	}
}
