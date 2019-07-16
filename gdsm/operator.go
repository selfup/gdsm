package gdsm

import (
	"bufio"
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

// recursive connection handler
func (op *Operator) handleConnection(conn net.Conn) {
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		op.handleReadConnErr(err, conn)
		conn.Close()
		return
	}

	message := strings.TrimSuffix(string(bufferBytes), "\n")
	clientAddr := conn.RemoteAddr()

	op.setNodes(clientAddr.String(), "")

	if !strings.Contains(message, " :: ") {
		if op.HandleSimplePayload(message, conn) {
			op.handleConnection(conn)
		} else {
			return
		}
	} else {
		op.HandleInstructionsPayload(message, conn)
		op.handleConnection(conn)
	}
}

func (op *Operator) registerServer(conn net.Conn, serverPort string) {
	client := conn.RemoteAddr().String()
	clientIP := strings.Split(client, ":")[0]
	server := clientIP + ":" + serverPort

	op.mutex.Lock()
	if !op.Servers[server] && server != "" {
		op.Servers[server] = true
	}

	servers := op.getServers()

	var wg sync.WaitGroup
	wg.Add(len(op.Servers))

	// update manager node client map to store servers
	op.Clients[client] = server

	serversString := strings.Join(servers, "|")

	for _, server := range servers {
		go func(serverAddr string) {
			Ping(serverAddr, "update_servers :: "+serversString)
			wg.Done()
		}(server)
	}

	wg.Wait()
	op.mutex.Unlock()
}

func (op *Operator) handleReadConnErr(err error, conn net.Conn) {
	op.removeConnFromCluster(conn)
}

func (op *Operator) removeConnFromCluster(conn net.Conn) {
	client := conn.RemoteAddr().String()
	op.deleteNode(client)

	op.mutex.Lock()
	servers := op.getServers()

	var wg sync.WaitGroup
	wg.Add(len(op.Clients))

	for key, value := range op.Clients {
		if value != "" {
			go func(clientAddr string, serverAddr string) {
				Ping(serverAddr, "remove_client :: "+client)
				Ping(serverAddr, "update_servers :: "+strings.Join(servers, "|"))
				wg.Done()
			}(key, value)
		} else {
			wg.Done()
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