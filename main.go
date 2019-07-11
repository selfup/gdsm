package main

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

const (
	netType = "tcp"
	netAddr = "127.0.0.1:8081"
)

var (
	cache = make(map[string]string)
	nodes = make(map[string]int)
	mutex = &sync.Mutex{}
)

func main() {
	listener, err := net.Listen(netType, netAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("JEANOME IS UP ON:", netAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}

func pingClient(conn net.Conn) {
	time.Sleep(100 * time.Millisecond)

	_, err := conn.Write([]byte("\n"))
	if err != nil {
		removeNodeFromCluster(conn)
	}

	pingClient(conn)
}

func handleConnection(conn net.Conn) {
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		handleReadConnErr(err, conn)
	}

	if err == io.EOF {
		handleReadConnErr(err, conn)
		conn.Close()
		return
	}

	bytes := bufferBytes
	message := strings.TrimSuffix(string(bufferBytes), "\n")
	clientAddr := conn.RemoteAddr()

	log.Println(message, bytes)

	checkForClientExistance(clientAddr)

	newMessage := strings.ToUpper(message)

	if !strings.Contains(newMessage, " :: ") {
		handleNonColonSpacedPayload(newMessage, conn)
		defer handleConnection(conn)
	} else {
		handleInstructionsPayload(newMessage, conn)
		defer handleConnection(conn)
	}
}

func handleNonColonSpacedPayload(newMessage string, conn net.Conn) {
	log.Println(newMessage)
	switch newMessage {
	case "EXIT":
		// remove IP addr from list locally and on all distributed nodes..
		deleteNode(conn.RemoteAddr().String())

		conn.Write([]byte("200\n"))
		conn.Close()
		break
	case "Q":
		conn.Write([]byte("200\n"))
		conn.Close()
		break
	default:
		conn.Write([]byte("405\n"))
		break
	}
}

func handleInstructionsPayload(newMessage string, conn net.Conn) {
	payload := strings.Split(newMessage, " :: ")
	verb := payload[0]

	switch verb {
	case "REGISTER":
		conn.Write([]byte("200\n"))
		break
	case "NODES":
		nodesStr := fmt.Sprintln(nodes)
		conn.Write([]byte(nodesStr + "\n"))
		break
	case "UNREGISTER":
		value := payload[1]
		deleteNode(value)
		conn.Write([]byte("\n"))
		break
	default:
		conn.Write([]byte("UNSUPPORTED INSTRUCTION\n"))
		break
	}
}

func handleReadConnErr(err error, conn net.Conn) {
	removeNodeFromCluster(conn)
	log.Println(conn.RemoteAddr(), "ReadConnErr:", err)
}

func removeNodeFromCluster(conn net.Conn) {
	value := conn.RemoteAddr().String()

	deleteNode(value)

	for key := range nodes {
		go func(ipAddr string) {

		}(key)
	}
}

func deleteNode(value string) {
	mutex.Lock()
	delete(nodes, value)
	mutex.Unlock()
}

func setNodes(value string) {
	mutex.Lock()
	nodes[value] = 1
	mutex.Unlock()
}

func checkForClientExistance(clientAddr net.Addr) {
	isAnExistingNode := false
	clientAddrStr := clientAddr.String()

	for key := range nodes {
		if key == clientAddrStr {
			isAnExistingNode = true
		}
	}

	if !isAnExistingNode {
		mutex.Lock()
		nodes[clientAddrStr] = 1
		mutex.Unlock()
	}

	log.Println(nodes)
}
