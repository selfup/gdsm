package main

import "github.com/selfup/jeanome"

func main() {
	operator := jeanome.New()
	operator.NetAddr = "127.0.0.1:8081"
	operator.Boot()
}
