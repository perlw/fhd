package main

import (
	"fmt"
	"math"
	"runtime"

	"github.com/perlw/fhd/internal/pkg/platform"
)

func setUp() {
	fmt.Println("SetUp")
}

func tearDown() {
	fmt.Println("TearDown")
}

var delta float64

func updateAndRender(backbuffer *platform.BitmapBuffer) {
	delta += 0.01
	zoom := 1 + math.Sin(delta*0.1)
	xmod, ymod := math.Sincos(delta)

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
		SetUp:           setUp,
		TearDown:        tearDown,
		UpdateAndRender: updateAndRender,
	}
	p.Main()
}
