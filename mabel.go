package main

import (
	"log"

	"github.com/eyedeekay/mabel/config/ini"
	"github.com/eyedeekay/mabel/tunnelmanager"
)

var tunnelManager tm.TunnelManager

func main() {
	sts, err := i2pini.SAMTunnelSlice("tunnels.ini")
	if err != nil {
		log.Fatal(err)
	}
	tm.NewTunnelManager("127.0.0.1", 7676, sts)
}
