package main

import (
	"github.com/selfup/gdsm/gdsm"
)

func main() {
	operator := gdsm.BuildGDSM()
	gdsm.BootMattDaemon(operator)
}
