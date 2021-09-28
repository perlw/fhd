package main

import (
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/perlw/fhd/internal/app/fhd/forza"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp4", "192.168.1.83:13337")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	server, err := net.ListenPacket("udp", ":0")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer server.Close()

	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	theta := 0.0
	for range ticker.C {
		theta += 0.1
		fVal := float32(math.Abs(math.Sin(theta)))

		var packet forza.DataPacket

		packet.Running = 1
		packet.CarPerformanceIndex = 999
		packet.CurrentEngineRpm = float32(rand.Intn(5000) + 800)
		packet.EngineIdleRpm = 800
		packet.EngineMaxRpm = 5800

		packet.Accel = 128
		packet.Brake = 128

		packet.NormalizedSuspensionFrontLeft = rand.Float32()
		packet.NormalizedSuspensionFrontRight = rand.Float32()
		packet.NormalizedSuspensionRearLeft = rand.Float32()
		packet.NormalizedSuspensionRearRight = rand.Float32()

		packet.TireTempFrontLeft = fVal * 200
		packet.TireTempFrontRight = fVal * 200
		packet.TireTempRearLeft = fVal * 200
		packet.TireTempRearRight = fVal * 200

		_, err := server.WriteTo(packet.ToBytes(), addr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
