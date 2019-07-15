package gdsm

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// Caller is the client struct for gdsm
type Caller struct {
	NetAddr string
	Server  string
}

// Dial reaches out to the Operator
func (c *Caller) Dial() {
	conn, err := net.Dial("tcp", c.NetAddr)
	if err != nil {
		log.Println(err, "..reconnecting")

		time.Sleep(1 * time.Second)

		c.Dial()
	}

	var serverAddress string

	if c.Server == "" {
		serverAddress = ""
	} else {
		serverAddress = c.Server
	}

	serverPort := strings.Split(serverAddress, ":")[1]

	size, err := fmt.Fprintf(conn, "register_server :: "+serverPort+"\n")
	if err != nil {
		log.Fatal(size, err)
	}

	connbuf := bufio.NewReader(conn)

	for {
		_, err := connbuf.ReadString('\n')

		if err != nil {
			log.Println(err, "..reconnecting")

			time.Sleep(1 * time.Second)

			c.Dial()
		}

		log.Println("dial tcp", conn.RemoteAddr().String()+":", "connect: ..connected")
	}
}

// Ping takes a remote server address and sends over a message from the host server
func Ping(netAddr string, message string) {
	conn, err := net.Dial("tcp", netAddr)
	if err != nil {
		log.Fatal(err)
	}

	size, err := fmt.Fprintf(conn, message+"\n")
	if err != nil {
		log.Fatal(size, err)
	}

	conn.Close()
}
