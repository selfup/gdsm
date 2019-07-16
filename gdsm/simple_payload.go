package gdsm

import (
	"fmt"
	"net"
)

// HandleSimplePayload  handles payloads without params
func (op *Operator) HandleSimplePayload(newMessage string, conn net.Conn) bool {
	switch newMessage {
	case "ping":
		conn.Write([]byte("200\n"))
		return true
	case "clients":
		nodesStr := fmt.Sprintln(op.Clients)
		conn.Write([]byte(nodesStr + "\n"))
		return true
	case "servers":
		op.mutex.Lock()
		servers := op.getServers()
		op.mutex.Unlock()

		serversStr := fmt.Sprintln(servers)
		conn.Write([]byte(serversStr + "\n"))
		return true
	case "workers":
		op.mutex.Lock()
		workers := op.getServers()
		op.mutex.Unlock()

		workersStr := fmt.Sprintln(workers)
		conn.Write([]byte(workersStr + "\n"))
		return true
	default:
		conn.Write([]byte("405\n"))
		return true
	}
}
