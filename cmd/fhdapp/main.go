package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/perlw/fhd/internal/pkg/forzaprotocol"
	"github.com/perlw/fhd/internal/pkg/platform"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

/*
   plotLine(int x0, int y0, int x1, int y1)
   dx =  abs(x1-x0);
   sx = x0<x1 ? 1 : -1;
   dy = -abs(y1-y0);
   sy = y0<y1 ? 1 : -1;
   err = dx+dy;  // error value e_xy
   while (true)   // loop
       plot(x0, y0);
       if (x0 == x1 && y0 == y1) break;
       e2 = 2*err;
       if (e2 >= dy) // e_xy+e_x > 0
           err += dy;
           x0 += sx;
       end if
       if (e2 <= dx) // e_xy+e_y < 0
           err += dx;
           y0 += sy;
       end if
   end while
*/
func clamp(a, min, max int) int {
	if a < min {
		return min
	} else if a > max {
		return max
	}
	return a
}

func drawLine(backbuffer *platform.BitmapBuffer, x1, y1, x2, y2 int) {
	var sx, sy, dx, dy int

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
		backbuffer.Memory[(y*backbuffer.Width)+x] = 0xffffffff
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
}

func (a *App) SetUp() {
	fmt.Println("SetUp")
}

func (a *App) TearDown() {
	fmt.Println("TearDown")
}

func (a *App) UpdateAndRender(backbuffer *platform.BitmapBuffer) {
	drawLine(backbuffer, 640, 0, 640, 720)
	drawLine(backbuffer, 0, 360, 1280, 360)
	drawLine(backbuffer, 0, 0, 1280, 720)
	drawLine(backbuffer, 0, 720, 1280, 0)

	x := (a.PosX / 5200) * 256
	y := -(a.PosZ / 5200) * 256
	backbuffer.Memory[(int(y+360)*backbuffer.Width)+int(x+640)] = 0xffff0000
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
		}
	}()

	p := platform.Platform{
		App: &app,
	}
	p.Main()
}
