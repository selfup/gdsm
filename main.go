package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

const (
	netType = "tcp"
	netAddr = "127.0.0.1:8081"
)

var (
	cache = make(map[string]string)
	nodes = make(map[string]string)
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

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		bytes := scanner.Bytes()
		message := scanner.Text()
		clientAddr := conn.RemoteAddr()

		log.Println(bytes)

		checkForClientExistance(clientAddr)

		newMessage := strings.ToUpper(message)

		if !strings.Contains(newMessage, " :: ") {
			handleNonColonSpacedPayload(newMessage, conn)
		} else {
			handleInstructionsPayload(newMessage, conn)
		}

		if err := scanner.Err(); err != nil {
			handleScannerErr(err)
		}
	}
}

func handleNonColonSpacedPayload(newMessage string, conn net.Conn) {
	switch newMessage {
	case "EXIT":
		// remove IP addr from list locally and on all distributed nodes..
		conn.Write([]byte("goodbye.." + "\n"))
		conn.Close()
		break
	case "Q":
		// remove IP addr from list locally and on all distributed nodes..
		conn.Write([]byte("goodbye.." + "\n"))
		conn.Close()
		break
	default:
		conn.Write([]byte("PAYLOAD NOT SUPPORTED\n"))
		break
	}
}

func removeNodeFromCluster(node string) {

}

func handleInstructionsPayload(newMessage string, conn net.Conn) {
	payload := strings.Split(newMessage, " :: ")
	verb := payload[0]
	key := payload[1]

	switch verb {
	case "GET":
		conn.Write([]byte(cache[key] + "\n"))
		break
	case "SET":
		value := payload[2]

		mutex.Lock()
		cache[key] = value
		mutex.Unlock()

		conn.Write([]byte("\n"))
		break
	case "DEL":
		value := payload[2]

		mutex.Lock()
		delete(cache, value)
		mutex.Unlock()

		break
	case "GET_NODES":
		nodesStr := fmt.Sprintln(nodes)
		conn.Write([]byte(nodesStr))
		break
	case "SET_NODES":
		value := payload[2]

		mutex.Lock()
		nodes[value] = "1"
		mutex.Unlock()

		break
	case "DEL_NODES":
		value := payload[2]

		mutex.Lock()
		delete(nodes, value)
		mutex.Unlock()

		break
	default:
		conn.Write([]byte("PAYLOAD NOT SUPPORTED\n"))
		break
	}
}

func handleScannerErr(err error) {
	switch strings.Contains(err.Error(), "use of closed network connection") {
	case true:
		// remove IP addr from list locally and on all distributed nodes..
		break
	default:
		log.Println("scanner.Err():", err)
		break
	}
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
		nodes[clientAddrStr] = "1"
		mutex.Unlock()
	}

	log.Println(nodes)
}
