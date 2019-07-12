package main

import "github.com/selfup/jeanome"

func main() {
	c := new(jeanome.Caller)

	c.NetAddr = "127.0.0.1:8081"

	c.Dial()
}
