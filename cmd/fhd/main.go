package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/perlw/fhd/internal/app/fhd/forza"
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
		class := []string{"D", "C", "B", "A", "S1", "S2", "R", "X"}
		gear := []rune{'R', '1', '2', '3', '4', '5', '6', '7'}
		for {
			_, _, err := listener.ReadFromUDP(buffer)
			if err != nil {
				fmt.Printf(err.Error())
			}

			var packet forza.DataPacket
			packet.FromBytes(buffer)
			if packet.Running == 0 {
				continue
			}

			//fmt.Printf("\rDAT: %X\n", buffer[:n])
			fmt.Printf(
				"\rID: %d (%s, PI: %d), Speed: %.2f, RPM: %.2f, Gear: %c",
				packet.CarID, class[packet.CarClass], packet.CarPerformanceIndex,
				packet.Speed*3.60, packet.CurrentEngineRpm, gear[packet.Gear],
			)
		}
	}()

	<-stop
}
