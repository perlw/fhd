package main

import (
	"fmt"
	"math"
	"runtime"

	"github.com/perlw/fhd/internal/pkg/platform"
)

type App struct {
	delta float64
}

func (a *App) SetUp() {
	fmt.Println("SetUp")
}

func (a *App) TearDown() {
	fmt.Println("TearDown")
}

func (a *App) UpdateAndRender(backbuffer *platform.BitmapBuffer) {
	a.delta += 0.01
	zoom := 1 + math.Sin(a.delta*0.1)
	xmod, ymod := math.Sincos(a.delta)

	for y := 0; y < backbuffer.Height; y++ {
		i := y * backbuffer.Width
		for x := 0; x < backbuffer.Width; x++ {
			cx := uint8((float64(x) + (xmod * 100)) * zoom)
			cy := uint8((float64(y) + (ymod * 100)) * zoom)
			m := uint32(cx ^ cy)
			c := (0xff << 24) | (m << 16) | (m << 8) | (m << 0)
			backbuffer.Memory[i+x] = c
		}
	}
}

func main() {
	runtime.LockOSThread()

	p := platform.Platform{
		App: &App{},
	}
	p.Main()
}
