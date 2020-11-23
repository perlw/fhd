package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/perlw/fhd/internal/app/fhd/forza"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:13337")
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	server, err := net.ListenPacket("udp", ":0")
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	defer server.Close()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		var packet forza.DataPacket
		packet.Running = 1
		packet.CurrentEngineRpm = float32(rand.Intn(5000) + 800)
		_, err := server.WriteTo(packet.ToBytes(), addr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
