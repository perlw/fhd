package forzaprotocol

import (
	"fmt"
	"net"
)

type Listener struct {
}

func (l *Listener) Listen(address string, dataChan chan<- Packet) error {
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return fmt.Errorf("could not resolve address: %w", err)
	}

	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("could not start listener: %w", err)
	}
	defer listener.Close()

	buffer := make([]byte, PacketSize)
	for {
		_, _, err := listener.ReadFromUDP(buffer)
		if err != nil {
			return fmt.Errorf("failed reading from socket: %w", err)
		}

		var p Packet
		p.FromBytes(buffer)
		dataChan <- p
	}
}
