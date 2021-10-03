package platform

type SetUpFunc func()
type TearDownFunc func()
type UpdateAndRenderFunc func(backbuffer *BitmapBuffer)

type BitmapBuffer struct {
	Memory []uint32
	Width  int
	Height int
	Bps    int
	Pitch  int
}

type Platform struct {
	SetUp           SetUpFunc
	TearDown        TearDownFunc
	UpdateAndRender UpdateAndRenderFunc
}
