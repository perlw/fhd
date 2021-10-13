package platform

type BitmapBuffer struct {
	Memory []uint32
	Width  int32
	Height int32
	Bps    int32
	Pitch  int32
}

type App interface {
	SetUp()
	TearDown()
	UpdateAndRender(backbuffer *BitmapBuffer)
}

type Platform struct {
	App App
}
