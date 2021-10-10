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
	var listenAddr string
	flag.StringVar(&listenAddr, "telemetry-addr", "", "which address:port to listen for telemetry data on")
	flag.Parse()

	if listenAddr == "" {
		fmt.Println("must supply a telemetry-addr!")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if _, err := net.ResolveUDPAddr("udp4", listenAddr); err != nil {
		fmt.Println(fmt.Errorf("invalid telemetry-addr: %w", err))
		os.Exit(1)
	}

	filename := fmt.Sprintf("%s.rec", time.Now().Format("200601021504"))
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	defer file.Close()

	var listener forzaprotocol.Listener
	dataChan := make(chan forzaprotocol.Packet)
	go func() {
		if err := listener.Listen(listenAddr, dataChan); err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
	}()

	var stopChan = make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	fmt.Printf("Recording packets from %s to file %s...\n", listenAddr, filename)

	var count int
quit:
	for {
		select {
		case packet := <-dataChan:
			c, _ := file.Write(packet.ToBytes())
			count += c
		case <-stopChan:
			break quit
		}
	}

	fmt.Printf("Recorded %.2fkbytes to %s!", float64(count)/1024, filename)
}
