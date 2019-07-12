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

	log.Println(size)

	if err != nil {
		log.Fatal(err)
	}

	connbuf := bufio.NewReader(conn)

	for {
		str, err := connbuf.ReadString('\n')

		if len(str) > 0 {
			log.Println(str)
		}

		if err != nil {
			log.Println(err)

			log.Println("attempting to reconnect..")

			time.Sleep(1 * time.Second)

			c.Dial()
		}
	}
}
