package main

import (
	"github.com/selfup/gdsm/gdsm"
)

func main() {
	daemon := gdsm.BuildDaemon()
	gdsm.BootDaemon(daemon)
}
