package main

import "github.com/selfup/jeanome"

func main() {
	caller := new(jeanome.Caller)
	caller.NetAddr = "127.0.0.1:8081"
	caller.Dial()
}
