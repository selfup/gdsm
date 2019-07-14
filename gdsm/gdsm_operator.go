package gdsm

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

// Operator is the TCP Server for gdsm
type Operator struct {
	mutex   sync.Mutex
	netType string
	NetAddr string
	Clients map[string]string
	Servers map[string]bool
}

// New is the Operator constructor
func New() *Operator {
	operator := new(Operator)

	operator.netType = "tcp"
	operator.NetAddr = "127.0.0.1:19888"
	operator.Clients = make(map[string]string)
	operator.Servers = make(map[string]bool)

	return operator
}

// Boot starts the gdsm TCP Server
func (op *Operator) Boot() {
	listener, err := net.Listen(op.netType, op.NetAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("GDSM IS UP ON:", op.NetAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go op.handleConnection(conn)
	}
}

func (op *Operator) handleConnection(conn net.Conn) {
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		op.handleReadConnErr(err, conn)
		conn.Close()
		return
	}

	message := strings.TrimSuffix(string(bufferBytes), "\n")
	clientAddr := conn.RemoteAddr()

	log.Println("IP", clientAddr.String(), "MESSAGE", message)

	op.setNodes(clientAddr.String(), "")

	if !strings.Contains(message, " :: ") {
		if op.handleSimplePayload(message, conn) {
			op.handleConnection(conn)
		} else {
			return
		}
	} else {
		op.handleInstructionsPayload(message, conn)
		op.handleConnection(conn)
	}
}

func (op *Operator) handleSimplePayload(newMessage string, conn net.Conn) bool {
	switch newMessage {
	case "register":
		conn.Write([]byte("200\n"))
		return true
	case "clients":
		nodesStr := fmt.Sprintln(op.Clients)
		conn.Write([]byte(nodesStr + "\n"))
		return true
	case "servers":
		op.mutex.Lock()
		var servers []string
		for server, active := range op.Servers {
			if active {
				servers = append(servers, server)
			}
		}
		op.mutex.Unlock()
		serversStr := fmt.Sprintln(servers)
		conn.Write([]byte(serversStr + "\n"))
		return true
	default:
		conn.Write([]byte("405\n"))
		return true
	}
}

func (op *Operator) registerServer(conn net.Conn, server string) {
	client := conn.RemoteAddr().String()

	op.mutex.Lock()
	if !op.Servers[server] && server != "" {
		op.Servers[server] = true
	}
	op.mutex.Unlock()

	servers := op.getServers()

	op.mutex.Lock()
	var wg sync.WaitGroup
	wg.Add(len(op.Servers))

	// update manager node client map to store servers
	op.Clients[client] = server

	for key, value := range op.Servers {
		if value && key != op.NetAddr {
			go func(serverAddr string) {
				Ping(serverAddr, "update_servers :: "+strings.Join(servers, "|"))
				wg.Done()
			}(key)
		} else {
			wg.Done()
		}
	}

	wg.Wait()
	op.mutex.Unlock()
}

func (op *Operator) handleInstructionsPayload(newMessage string, conn net.Conn) {
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

func (op *Operator) handleReadConnErr(err error, conn net.Conn) {
	if err == io.EOF {
		log.Println("IP", conn.RemoteAddr(), "disconnected..")
	} else {
		log.Println("IP", conn.RemoteAddr(), "ERR", err)
	}

	op.removeConnFromCluster(conn)
}

func (op *Operator) removeConnFromCluster(conn net.Conn) {
	client := conn.RemoteAddr().String()

	op.deleteNode(client)
	servers := op.getServers()

	op.mutex.Lock()
	var wg sync.WaitGroup
	wg.Add(len(op.Clients))

	for key, value := range op.Clients {
		if value != "" && value != op.NetAddr {
			go func(clientAddr string, serverAddr string) {
				log.Println("REMOVE_CONN_FROM_CLUSTER CALLING", serverAddr, "REMOVING CLIENT", client)
				Ping(serverAddr, "remove_client :: "+client)
				Ping(serverAddr, "update_servers :: "+strings.Join(servers, "|"))
				wg.Done()
			}(key, value)
		}
	}

	wg.Wait()
	op.mutex.Unlock()
}

func (op *Operator) deleteNode(value string) {
	op.mutex.Lock()
	// when a stored client has an associated server then remove the server
	if op.Clients[value] != "" {
		delete(op.Servers, op.Clients[value])
	}

	delete(op.Clients, value)
	op.mutex.Unlock()
}

func (op *Operator) setNodes(key string, value string) {
	op.mutex.Lock()
	op.Clients[key] = value
	op.mutex.Unlock()
}

func (op *Operator) getServers() []string {
	var servers []string

	for server, active := range op.Servers {
		if active {
			servers = append(servers, server)
		}
	}

	return servers
}

func (op *Operator) updateServers(value string) {
	servers := strings.Split(value, "|")

	op.mutex.Lock()
	for serv := range op.Servers {
		op.Servers[serv] = false
	}

	for _, server := range servers {
		op.Servers[server] = true
	}
	op.mutex.Unlock()
}
