package jeanome

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

// Caller is the client struct for jeanome
type Caller struct {
	NetAddr string
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

	size, err := fmt.Fprintf(conn, "register\n")
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
