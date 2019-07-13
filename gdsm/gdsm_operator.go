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
	Nodes   map[string]string
}

// New is the Operator constructor
func New() *Operator {
	operator := new(Operator)

	operator.netType = "tcp"
	operator.NetAddr = "127.0.0.1:19888"
	operator.Nodes = make(map[string]string)

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

	if err == io.EOF {
		op.handleReadConnErr(err, conn)
		conn.Close()
		return
	}

	if err != nil {
		op.handleReadConnErr(err, conn)
	}

	message := strings.TrimSuffix(string(bufferBytes), "\n")
	clientAddr := conn.RemoteAddr()

	log.Println("IP", clientAddr.String(), "MESSAGE", message)

	op.setNodes(clientAddr.String(), "..")

	newMessage := strings.ToUpper(message)

	if !strings.Contains(newMessage, " :: ") {
		if op.handleSimplePayload(newMessage, conn) {
			op.handleConnection(conn)
		} else {
			return
		}
	} else {
		op.handleInstructionsPayload(newMessage, conn)
		op.handleConnection(conn)
	}
}

func (op *Operator) handleSimplePayload(newMessage string, conn net.Conn) bool {
	switch newMessage {
	case "REGISTER":
		conn.Write([]byte("200\n"))
		return true
	case "NODES":
		nodesStr := fmt.Sprintln(op.Nodes)
		conn.Write([]byte(nodesStr + "\n"))
		return true
	default:
		conn.Write([]byte("405\n"))
		return true
	}
}

func (op *Operator) registerServer(conn net.Conn, server string) {
	client := conn.RemoteAddr().String()

	op.mutex.Lock()
	op.Nodes[client] = server
	op.mutex.Unlock()
}

func (op *Operator) handleInstructionsPayload(newMessage string, conn net.Conn) {
	payload := strings.Split(newMessage, " :: ")
	verb := payload[0]

	switch verb {
	case "REMOVE_CLIENT":
		value := payload[1]

		op.removeNodeFromCluster(value)

		conn.Close()
		conn = nil

		break
	case "REGISTER_SERVER":
		value := payload[1]

		op.registerServer(conn, value)
		log.Println(op.Nodes)

		conn.Write([]byte("200\n"))
		break
	case "UNREGISTER":
		value := payload[1]

		op.removeNodeFromCluster(value)

		conn.Write([]byte("200\n"))
		break
	default:
		conn.Write([]byte("UNSUPPORTED INSTRUCTION\n"))
		break
	}
}

func (op *Operator) handleReadConnErr(err error, conn net.Conn) {
	log.Println("IP", conn.RemoteAddr(), "ERR", err)
	op.removeConnFromCluster(conn)
}

func (op *Operator) removeNodeFromCluster(node string) {
	op.deleteNode(node)

	var wg sync.WaitGroup

	wg.Add(len(op.Nodes))

	op.mutex.Lock()
	for key, value := range op.Nodes {
		if value != ".." && value != op.NetAddr {
			go func(clientAddr string, serverAddr string) {
				log.Println("REMOVE_NODE_FROM_CLUSTER CALLING", serverAddr, "REMOVING CLIENT", clientAddr)
				Ping(serverAddr, "remove_client ::"+clientAddr)
				wg.Done()
			}(key, value)
		}
	}
	op.mutex.Unlock()

	wg.Wait()
}

func (op *Operator) removeConnFromCluster(conn net.Conn) {
	client := conn.RemoteAddr().String()

	op.deleteNode(client)

	var wg sync.WaitGroup

	wg.Add(len(op.Nodes))

	op.mutex.Lock()
	for key, value := range op.Nodes {
		if value != ".." && value != op.NetAddr {
			go func(clientAddr string, serverAddr string) {
				log.Println("REMOVE_CONN_FROM_CLUSTER CALLING", serverAddr, "REMOVING CLIENT", client)
				Ping(serverAddr, "remove_client :: "+client)
				wg.Done()
			}(key, value)
		}
	}
	op.mutex.Unlock()

	wg.Wait()
}

func (op *Operator) deleteNode(value string) {
	op.mutex.Lock()
	delete(op.Nodes, value)
	op.mutex.Unlock()
}

func (op *Operator) setNodes(key string, value string) {
	op.mutex.Lock()
	op.Nodes[key] = value
	op.mutex.Unlock()
}
