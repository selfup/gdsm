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

func pingClient(conn net.Conn) {
	time.Sleep(100 * time.Millisecond)

	_, err := conn.Write([]byte("\n"))
	if err != nil {
		removeNodeFromCluster(conn.RemoteAddr().String())
	}

	pingClient(conn)
}

func handleConnection(conn net.Conn) {
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		handleScannerErr(err)
		log.Println(err)
	}

	if err == io.EOF {
		handleScannerErr(err)
		conn.Close()
		log.Println(err)
		return
	}

	bytes := bufferBytes
	message := strings.TrimSuffix(string(bufferBytes), "\n")
	clientAddr := conn.RemoteAddr()

	log.Println(message, bytes)

	checkForClientExistance(clientAddr)

	newMessage := strings.ToUpper(message)

	if !strings.Contains(newMessage, " :: ") {
		log.Println("Non Colon")
		handleNonColonSpacedPayload(newMessage, conn)
		defer handleConnection(conn)
	} else {
		log.Println("With Colon")
		handleInstructionsPayload(newMessage, conn)
		defer handleConnection(conn)
	}
}

func handleNonColonSpacedPayload(newMessage string, conn net.Conn) {
	log.Println(newMessage)
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
		log.Println(cache[key])
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
	// remove IP addr from list locally and on all distributed nodes..
	log.Println("scanner.Err():", err)
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
