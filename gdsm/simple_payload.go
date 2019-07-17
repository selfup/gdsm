package gdsm

import (
	"fmt"
	"net"
)

// HandleSimplePayload  handles payloads without params
func (op *Operator) HandleSimplePayload(newMessage string, conn net.Conn) {
	switch newMessage {
	case "ping":
		conn.Write([]byte("200\n"))
		break
	case "clients":
		nodesStr := fmt.Sprintln(op.Clients)
		conn.Write([]byte(nodesStr + "\n"))
		break
	case "servers":
		op.mutex.Lock()
		servers := op.getServers()
		op.mutex.Unlock()

		serversStr := fmt.Sprintln(servers)
		conn.Write([]byte(serversStr + "\n"))
		break
	case "workers":
		op.mutex.Lock()
		workers := op.getServers()
		op.mutex.Unlock()

		workersStr := fmt.Sprintln(workers)
		conn.Write([]byte(workersStr + "\n"))
		break
	default:
		conn.Write([]byte("405\n"))
		break
	}
}
