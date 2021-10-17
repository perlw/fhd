package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/perlw/fhd/internal/pkg/forzaprotocol"
	"github.com/perlw/fhd/internal/pkg/platform"
)

func abs(a int32) int32 {
	v := a >> 31
	return (a + v) ^ v
}

func clamp(a, min, max int32) int32 {
	if a < min {
		return min
	} else if a > max {
		return max
	}
	return a
}

func lerp(start, stop int32, t float32) int32 {
	return int32(((1 - t) * float32(start)) + (t * float32(stop)) + 0.5)
}

func drawLine(backbuffer *platform.BitmapBuffer, x1, y1, x2, y2 int32, color uint32) {
	var sx, sy, dx, dy int32

	x1 = clamp(x1, 0, backbuffer.Width-1)
	y1 = clamp(y1, 0, backbuffer.Height-1)
	x2 = clamp(x2, 0, backbuffer.Width-1)
	y2 = clamp(y2, 0, backbuffer.Height-1)

	dx = abs(x2 - x1)
	if x1 < x2 {
		sx = 1
	} else {
		sx = -1
	}
	dy = -abs(y2 - y1)
	if y1 < y2 {
		sy = 1
	} else {
		sy = -1
	}

	x := x1
	y := y1
	e := dx + dy
	for {
		backbuffer.Memory[(y*backbuffer.Width)+x] = color
		if x == x2 && y == y2 {
			break
		}

		e2 := e * 2
		if e2 >= dy {
			e += dy
			x += sx
		}
		if e2 <= dx {
			e += dx
			y += sy
		}
	}
}

type App struct {
	PosX, PosZ float32
	Gas, Brake float32
	Laps       uint16
	CurrentLap uint16
}

type State struct {
	AccelMinShade, AccelMaxShade int32
	BrakeMinShade, BrakeMaxShade int32
}

func (a *App) SetUp(memory *platform.Memory) {
	state := (*State)(memory.PermanentStorage)
	state.AccelMinShade = 96
	state.AccelMaxShade = 255
	state.BrakeMinShade = 96
	state.BrakeMaxShade = 255
}

func (a *App) TearDown() {
	fmt.Println("Bye")
}

func (a *App) UpdateAndRender(memory *platform.Memory, backbuffer *platform.BitmapBuffer, elapsedMs float64) {
	fmt.Printf("elapsed: %f\n", elapsedMs)

	state := (*State)(memory.PermanentStorage)

	if a.Laps != a.CurrentLap {
		a.Laps = a.CurrentLap
		for i, c := range backbuffer.Memory {
			r := ((c >> 16) & 0xff) >> 2
			g := ((c >> 8) & 0xff) >> 2
			b := ((c >> 0) & 0xff) >> 2
			backbuffer.Memory[i] = (0xff << 24) + (r << 16) + (g << 8) + (b << 0)
		}
	}

	drawLine(backbuffer, 640, 0, 640, 720, 0xff333333)
	drawLine(backbuffer, 0, 360, 1280, 360, 0xff333333)
	drawLine(backbuffer, 0, 0, 1280, 720, 0xff333333)
	drawLine(backbuffer, 0, 720, 1280, 0, 0xff333333)

	drawLine(backbuffer, 384, 104, 896, 104, 0xffffffff)
	drawLine(backbuffer, 896, 104, 896, 616, 0xffffffff)
	drawLine(backbuffer, 896, 616, 384, 616, 0xffffffff)
	drawLine(backbuffer, 384, 104, 384, 616, 0xffffffff)

	var viewTLX, viewTLZ float32 = -1000, -5300
	var viewBRX, viewBRZ float32 = 200, -6500
	viewWidth, viewHeight := viewBRX-viewTLX, viewTLZ-viewBRZ
	modPosX, modPosZ := (a.PosX-viewTLX)/float32(viewWidth), (a.PosZ-viewBRZ)/float32(viewHeight)

	if modPosX >= 0 && modPosX <= 1 && modPosZ >= 0 && modPosZ <= 1 {
		x := ((modPosX * 2) - 1) * 256
		y := ((-modPosZ * 2) + 1) * 256

		r := uint32(lerp(state.AccelMinShade, state.AccelMaxShade, a.Brake))
		g := uint32(lerp(state.BrakeMinShade, state.BrakeMaxShade, a.Gas))
		b := uint32(96)

		backbuffer.Memory[(int32(y+360)*backbuffer.Width)+int32(x+640)] = (0xff << 24) + (r << 16) + (g << 8) + (b << 0)
	}
}

func main() {
	runtime.LockOSThread()

	var app App

	dataChan := make(chan forzaprotocol.Packet)
	var listener forzaprotocol.Listener
	go func() {
		if err := listener.Listen("0.0.0.0:13337", dataChan); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}()
	go func() {
		for packet := range dataChan {
			app.PosX = packet.PositionX
			app.PosZ = packet.PositionZ
			app.Gas = float32(packet.Accel) / 255
			app.Brake = float32(packet.Brake) / 255
			app.CurrentLap = packet.LapNumber
		}
	}()

	p := platform.Platform{
		App: &app,
	}
	p.Main()
}
