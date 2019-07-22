package main

import (
	"github.com/selfup/gdsm"
)

func main() {
	daemon := gdsm.BuildDaemon()
	gdsm.BootDaemon(daemon)
}
