package jeanome

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// Operator is the TCP Server for Jeanome
type Operator struct {
	mutex   sync.Mutex
	netType string
	cache   map[string]string
	NetAddr string
	Nodes   map[string]int
}

// New is the Operator constructor
func New() *Operator {
	operator := new(Operator)

	operator.netType = "tcp"
	operator.cache = make(map[string]string)
	operator.NetAddr = "127.0.0.1:19888"
	operator.Nodes = make(map[string]int)

	return operator
}

// Boot starts the Jeanome TCP Server
func (op *Operator) Boot() {
	listener, err := net.Listen(op.netType, op.NetAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("JEANOME IS UP ON:", op.NetAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go op.handleConnection(conn)
	}
}

func (op *Operator) pingClient(conn net.Conn) {
	time.Sleep(100 * time.Millisecond)

	_, err := conn.Write([]byte("\n"))
	if err != nil {
		op.removeConnFromCluster(conn)
	}

	op.pingClient(conn)
}

func (op *Operator) handleConnection(conn net.Conn) {
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		op.handleReadConnErr(err, conn)
	}

	if err == io.EOF {
		op.handleReadConnErr(err, conn)
		conn.Close()
		return
	}

	bytes := bufferBytes
	message := strings.TrimSuffix(string(bufferBytes), "\n")
	clientAddr := conn.RemoteAddr()

	log.Println(message, bytes)

	op.mutex.Lock()
	op.checkForClientExistance(clientAddr)
	op.mutex.Unlock()

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
	case "EXIT":
		// remove IP addr from list locally and on all distributed nodes..
		op.deleteNode(conn.RemoteAddr().String())

		conn.Close()
		conn = nil

		return false
	case "Q":
		conn.Close()
		conn = nil

		return false
	default:
		conn.Write([]byte("405\n"))
		return true
	}
}

func (op *Operator) handleInstructionsPayload(newMessage string, conn net.Conn) {
	payload := strings.Split(newMessage, " :: ")
	verb := payload[0]

	switch verb {
	case "UNREGISTER":
		value := payload[1]

		op.deleteNode(value)
		op.removeNodeFromCluster(value)

		conn.Write([]byte("200\n"))
		break
	default:
		conn.Write([]byte("UNSUPPORTED INSTRUCTION\n"))
		break
	}
}

func (op *Operator) handleReadConnErr(err error, conn net.Conn) {
	log.Println(conn.RemoteAddr(), "ReadConnErr:", err)
	op.removeConnFromCluster(conn)
}

func (op *Operator) removeNodeFromCluster(node string) {
	value := net.ParseIP(node)

	op.deleteNode(value.String())

	for key := range op.Nodes {
		go func(ipAddr string) {

		}(key)
	}
}

func (op *Operator) removeConnFromCluster(conn net.Conn) {
	value := conn.RemoteAddr().String()

	op.deleteNode(value)

	op.mutex.Lock()
	for key := range op.Nodes {
		go func(ipAddr string) {

		}(key)
	}
	op.mutex.Unlock()
}

func (op *Operator) deleteNode(value string) {
	op.mutex.Lock()
	delete(op.Nodes, value)
	op.mutex.Unlock()
}

func (op *Operator) setNodes(value string) {
	op.mutex.Lock()
	op.Nodes[value] = 1
	op.mutex.Unlock()
}

func (op *Operator) checkForClientExistance(clientAddr net.Addr) {
	isAnExistingNode := false
	clientAddrStr := clientAddr.String()

	for key := range op.Nodes {
		if key != "" && key == clientAddrStr {
			isAnExistingNode = true
		}
	}

	if !isAnExistingNode {
		op.Nodes[clientAddrStr] = 1
	}

	log.Println(op.Nodes)
}
