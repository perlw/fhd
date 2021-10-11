package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/perlw/fhd/internal/pkg/forzaprotocol"
)

func main() {
	if err := termui.Init(); err != nil {
		fmt.Printf("failed to initialize termui: %v", err)
	}
	defer termui.Close()

	infoBar := widgets.NewParagraph()
	infoBar.Text = "WAITING FOR DATA"
	infoBar.SetRect(0, 0, 80, 3)
	rpmGauge := widgets.NewGauge()
	rpmGauge.Title = "RPM"
	rpmGauge.SetRect(0, 3, 20, 6)
	positionBar := widgets.NewParagraph()
	positionBar.Title = "Position"
	positionBar.SetRect(21, 3, 60, 6)
	suspension := widgets.NewBarChart()
	suspension.Title = "Suspension"
	suspension.NumFormatter = func(a float64) string {
		return fmt.Sprintf("%0.1f", a)
	}
	suspension.SetRect(0, 6, 18, 14)
	suspension.Data = []float64{0.5, 0.5, 0.5, 0.5}
	suspension.Labels = []string{"FL", "FR", "RL", "RR"}
	suspension.MaxVal = 1.0
	temperature := widgets.NewBarChart()
	temperature.Title = "Tire Temp"
	temperature.SetRect(19, 6, 37, 14)
	temperature.NumFormatter = func(a float64) string {
		return fmt.Sprintf("%0.0f", a)
	}
	temperature.Data = []float64{0.5, 0.5, 0.5, 0.5}
	temperature.Labels = []string{"FL", "FR", "RL", "RR"}
	temperature.MaxVal = 200.0
	control := widgets.NewBarChart()
	control.Title = "Control"
	control.NumFormatter = func(a float64) string {
		return fmt.Sprintf("%d", uint8(a))
	}
	control.SetRect(38, 6, 48, 14)
	control.Data = []float64{0, 0}
	control.Labels = []string{"Acc", "Brake"}
	control.MaxVal = 255

	termui.Render(infoBar)

	class := []string{"D", "C", "B", "A", "S1", "S2", "R", "X"}
	gear := []rune{'R', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'A', 'B', 'C', 'D'}

	packet := make(chan forzaprotocol.Packet)
	var listener forzaprotocol.Listener
	go func() {
		if err := listener.Listen("0.0.0.0:13337", packet); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}()

	uiEvents := termui.PollEvents()
	var stopChan = make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

quit:
	for {
		select {
		case p := <-packet:
			if p.Running == 0 {
				break
			}

			infoBar.Text = fmt.Sprintf(
				"ID: %d (%s, PI: %d), Speed: %.2f, RPM: %.2f, Gear: %c",
				p.CarID, class[p.CarClass], p.CarPerformanceIndex,
				p.Speed*3.60, p.CurrentEngineRpm, gear[p.Gear],
			)

			/**
			 * FH4 - House positions tested
			 *                X    Y    Z
			 * Top Left     -3704 160  1343
			 * Top Center   -1455 152  3552
			 * Top Right     2801 166  1499
			 * Bottom Left  -5191 156 -3517
			 * Bottom Center -242 227 -5565
			 * Bottom Right  5208 115 -3388
			 *
			 * Guesstimated map boundaries:
			 *       X     Z
			 * TL: -5191  3552
			 * BR:  5208 -5565
			 */
			positionBar.Text = fmt.Sprintf(
				"X: %.2f, Y: %.2f, Z: %.2f",
				p.PositionX, p.PositionY, p.PositionZ,
			)

			rpmGauge.Percent = int((p.CurrentEngineRpm / p.EngineMaxRpm) * 100.0)
			suspension.Data[0] = float64(p.NormalizedSuspensionFrontLeft)
			suspension.Data[1] = float64(p.NormalizedSuspensionFrontRight)
			suspension.Data[2] = float64(p.NormalizedSuspensionRearLeft)
			suspension.Data[3] = float64(p.NormalizedSuspensionRearRight)
			temperature.Data[0] = float64(p.TireTempFrontLeft)
			temperature.Data[1] = float64(p.TireTempFrontRight)
			temperature.Data[2] = float64(p.TireTempRearLeft)
			temperature.Data[3] = float64(p.TireTempRearRight)
			control.Data[0] = float64(p.Accel)
			control.Data[1] = float64(p.Brake)

			termui.Render(infoBar)
			termui.Render(rpmGauge)
			termui.Render(positionBar)
			termui.Render(suspension)
			termui.Render(temperature)
			termui.Render(control)

		case e := <-uiEvents:
			if e.Type == termui.KeyboardEvent {
				break quit
			}

		case <-stopChan:
			break quit
		}
	}

	fmt.Println("Bye!")
}
