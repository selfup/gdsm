package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	pingMaxFiles()
}

func pingMaxFiles() {
	var wg sync.WaitGroup

	arg := os.Args[1]
	requests := 1024

	wg.Add(requests)

	for i := 0; i < requests; i++ {
		go func() {
			call(arg)
			wg.Done()
		}()
	}

	wg.Wait()
}

func call(netAddr string) {
	conn, err := net.Dial("tcp", netAddr)
	if err != nil {
		log.Fatal(err)
	}

	_, ferr := fmt.Fprintf(conn, "ping\n")
	if ferr != nil {
		log.Fatal(ferr)
	}

	conn.Close()
}
