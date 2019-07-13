package main

import "github.com/selfup/gdsm/gdsm"

func main() {
	operator := gdsm.New()
	operator.NetAddr = "127.0.0.1:8081"
	operator.Boot()
}
