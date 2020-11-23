package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/perlw/fhd/internal/app/fhd/forza"
)

func listen(packet chan<- forza.DataPacket) {
	addr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:13337")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	buffer := make([]byte, 1500)
	for {
		_, _, err := listener.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err.Error())
		}

		var dp forza.DataPacket
		dp.FromBytes(buffer)
		if dp.Running == 0 {
			continue
		}

		packet <- dp
	}
}

func main() {
	if err := termui.Init(); err != nil {
		fmt.Printf("failed to initialize termui: %v", err)
	}
	defer termui.Close()

	pg := widgets.NewParagraph()
	pg.Text = "ID: 0 (-, PI: 0), Speed: 0.00, RPM: 0.00, Gear: -"
	pg.SetRect(0, 0, 75, 3)
	g := widgets.NewGauge()
	g.SetRect(0, 3, 75, 6)

	termui.Render(pg)
	termui.Render(g)

	class := []string{"D", "C", "B", "A", "S1", "S2", "R", "X"}
	gear := []rune{'R', '1', '2', '3', '4', '5', '6', '7'}
	packet := make(chan forza.DataPacket)
	go listen(packet)

	var end bool
	uiEvents := termui.PollEvents()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	for !end {
		select {
		case p := <-packet:
			text := fmt.Sprintf(
				"ID: %d (%s, PI: %d), Speed: %.2f, RPM: %.2f, Gear: %c",
				p.CarID, class[p.CarClass], p.CarPerformanceIndex,
				p.Speed*3.60, p.CurrentEngineRpm, gear[p.Gear],
			)

			pg.Text = text
			g.Percent = int((p.CurrentEngineRpm / p.EngineMaxRpm) * 100.0)

			termui.Render(pg)
			termui.Render(g)
		case e := <-uiEvents:
			if e.Type == termui.KeyboardEvent {
				end = true
				break
			}
		case <-stop:
			end = true
			break
		default:
		}
	}
}
