package gdsm

import (
	"net"
	"strings"
)

// Update handles complex actions that involve params
func (op *Operator) Update(newMessage string, conn net.Conn) {
	payload := strings.Split(newMessage, " :: ")
	verb := payload[0]

	switch verb {
	case "remove_client":
		value := payload[1]
		op.deleteNode(value)
		conn.Write([]byte("200\n"))
		break
	case "register_server":
		value := payload[1]
		op.registerServer(conn, value)
		conn.Write([]byte("200\n"))
		break
	case "update_servers":
		value := payload[1]
		op.updateServers(value)
		conn.Write([]byte("200\n"))
		break
	default:
		conn.Write([]byte("405\n"))
		break
	}
}
