package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")

	if ip == "" {
		ip = "127.0.0.1"
	}

	if port == "" {
		port = "8081"
	}

	addr := ip + ":" + port

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		reader := bufio.NewReader(os.Stdin)

		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(conn, text+"\n")

		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(message)
	}
}
