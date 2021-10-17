package platform

import "unsafe"

type BitmapBuffer struct {
	Memory []uint32
	Width  int32
	Height int32
	Bps    int32
	Pitch  int32
}

type Memory struct {
	PermanentSize    uint64
	PermanentStorage unsafe.Pointer
}

type App interface {
	SetUp(memory *Memory)
	TearDown()
	UpdateAndRender(memory *Memory, backbuffer *BitmapBuffer, elapsedMs float64)
}

type Platform struct {
	App App
}
