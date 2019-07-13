package gdsm

import (
	"bufio"
	"fmt"
	"log"
	"net"
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
		log.Println(err)

		log.Println("attempting to reconnect..")

		time.Sleep(1 * time.Second)

		c.Dial()
	}

	var serverAddress string

	if c.Server == "" {
		serverAddress = ".."
	} else {
		serverAddress = c.Server
	}

	size, err := fmt.Fprintf(conn, "register_server :: "+serverAddress+"\n")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("IP", conn.RemoteAddr().String(), "BYTES", size, "connected..")

	connbuf := bufio.NewReader(conn)

	for {
		_, err := connbuf.ReadString('\n')

		if err != nil {
			log.Println(err)

			log.Println("attempting to reconnect..")

			time.Sleep(1 * time.Second)

			c.Dial()
		}
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
		log.Fatal(err)
	}

	log.Println("IP", netAddr, "SIZE", size)

	conn.Close()
}
