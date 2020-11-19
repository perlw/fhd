package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:13337")
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("listening...\n")
		buffer := make([]byte, 1500)
		for {
			n, _, err := listener.ReadFromUDP(buffer)
			if err != nil {
				fmt.Printf(err.Error())
			}

			fmt.Printf("DAT: %X\n\n", buffer[:n])
		}
	}()

	<-stop
}
