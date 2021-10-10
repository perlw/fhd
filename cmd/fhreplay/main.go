package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/perlw/fhd/internal/pkg/forzaprotocol"
)

func main() {
	var telemetryAddr string
	var filename string

	flag.StringVar(&filename, "recording", "", "which recording to playback")
	flag.StringVar(&telemetryAddr, "telemetry-addr", "", "which address:port to put telemetry data on")
	flag.Parse()

	if filename == "" {
		fmt.Println("must supply a recording filename!")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if telemetryAddr == "" {
		fmt.Println("must supply a telemetry-addr!")
		flag.PrintDefaults()
		os.Exit(1)
	}
	addr, err := net.ResolveUDPAddr("udp4", telemetryAddr)
	if err != nil {
		fmt.Println(fmt.Errorf("invalid telemetry-addr: %w", err))
		os.Exit(1)
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	defer file.Close()

	server, err := net.ListenPacket("udp", ":0")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer server.Close()

	ticker := time.NewTicker(time.Millisecond * 10)
	defer ticker.Stop()

	var stopChan = make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	fmt.Printf("Playing back packets from file %s to %s...\n", filename, telemetryAddr)

quit:
	for {
		select {
		case <-ticker.C:
			buffer := make([]byte, forzaprotocol.PacketSize)
			n, _ := file.Read(buffer)
			if n == 0 {
				break quit
			}
			if n != int(forzaprotocol.PacketSize) {
				fmt.Printf("Something went wrong, expected to read %d bytes but read %d\n", forzaprotocol.PacketSize, n)
				os.Exit(2)
			}
			_, err := server.WriteTo(buffer, addr)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(2)
			}
		case <-stopChan:
			break quit
		}
	}

	fmt.Printf("File %s played back successfully to %s\n", filename, telemetryAddr)
}
