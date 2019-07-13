package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func call(netAddr string) {
	conn, err := net.Dial("tcp", netAddr)
	if err != nil {
		log.Fatal(err)
	}

	size, err := fmt.Fprintf(conn, "register\n")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(size)

	conn.Close()
}

func main() {
	for {
		call(os.Args[1])
		time.Sleep(100 * time.Millisecond)
	}
}
