package main

import "github.com/selfup/gdsm/gdsm"

func main() {
	caller := new(gdsm.Caller)
	caller.NetAddr = "127.0.0.1:8081"
	caller.Dial()
}
