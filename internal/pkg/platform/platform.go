package platform

type BitmapBuffer struct {
	Memory []uint32
	Width  int
	Height int
	Bps    int
	Pitch  int
}

type App interface {
	SetUp()
	TearDown()
	UpdateAndRender(backbuffer *BitmapBuffer)
}

type Platform struct {
	App App
}
