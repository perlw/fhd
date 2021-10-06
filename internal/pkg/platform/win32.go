//go:build windows

package platform

import (
	"fmt"
	"reflect"
	"unsafe"
)

/*
#cgo CFLAGS:-std=c99
#cgo LDFLAGS: -luser32 -lgdi32
#include "win32.h"
*/
import "C"

type backbufferInfo struct {
	bitmapInfo C.BITMAPINFO
	memory     unsafe.Pointer
	width      int
	height     int
	bps        int
	pitch      int
}

var globalIsRunning bool
var globalBackbuffer backbufferInfo

func resizeDIBSection(backbuffer *backbufferInfo, width, height int) {
	if backbuffer.memory != nil {
		C.VirtualFree(C.LPVOID(backbuffer.memory), 0, C.MEM_RELEASE)
	}

	backbuffer.width = width
	backbuffer.height = height

	backbuffer.bitmapInfo.bmiHeader = C.BITMAPINFOHEADER{
		biSize:        C.ulong(unsafe.Sizeof(backbuffer.bitmapInfo.bmiHeader)),
		biWidth:       C.long(backbuffer.width),
		biHeight:      C.long(-backbuffer.height),
		biPlanes:      1,
		biBitCount:    32,
		biCompression: C.BI_RGB,
	}

	backbuffer.bps = 4
	backbuffer.pitch = backbuffer.width * backbuffer.bps
	memorySize := C.ulonglong(backbuffer.bps * (backbuffer.width * backbuffer.height))
	backbuffer.memory = unsafe.Pointer(C.VirtualAlloc(nil, memorySize, C.MEM_RESERVE|C.MEM_COMMIT, C.PAGE_READWRITE))
}

func blitBufferInWindow(backbuffer *backbufferInfo, dc C.HDC, width, height int) {
	var check float32 = 16.0 / 9.0
	correctedWidth := int(float32(height) * check)
	offsetX := (width - correctedWidth) / 2
	if correctedWidth != width {
		C.PatBlt(dc, 0, 0, C.int(offsetX), C.int(height), C.BLACKNESS)
		C.PatBlt(dc, C.int(width-offsetX), 0, C.int(offsetX), C.int(height), C.BLACKNESS)
	}
	C.StretchDIBits(dc, C.int(offsetX), 0, C.int(correctedWidth), C.int(height), 0, 0,
		C.int(backbuffer.width), C.int(backbuffer.height), backbuffer.memory,
		&backbuffer.bitmapInfo, C.DIB_RGB_COLORS, C.SRCCOPY)
}

//export WindowProc
func WindowProc(window C.HWND, message C.UINT, wParam C.WPARAM, lParam C.LPARAM) C.LRESULT {
	var result C.LRESULT = 0

	switch message {
	case C.WM_DESTROY:
		fmt.Println("WM_DESTROY")
		globalIsRunning = false

	case C.WM_CLOSE:
		fmt.Println("WM_CLOSE")
		globalIsRunning = false

	case C.WM_ACTIVATEAPP:
		fmt.Println("WM_ACTIVATEAPP")

	case C.WM_PAINT:
		fmt.Println("WM_PAINT")

		var paint C.PAINTSTRUCT
		dc := C.BeginPaint(window, &paint)
		var clientRect C.RECT
		C.GetClientRect(window, &clientRect)
		blitBufferInWindow(&globalBackbuffer, dc, int(clientRect.right-clientRect.left), int(clientRect.bottom-clientRect.top))
		C.EndPaint(window, &paint)

	default:
		result = C.DefWindowProc(window, message, wParam, lParam)
	}

	return result
}

func (p *Platform) Main() {
	className := C.CString("fhdwin32platform")
	hInstance := C.GetModuleHandle(nil)
	windowClass := C.WNDCLASSA{
		style:         C.CS_OWNDC | C.CS_HREDRAW | C.CS_VREDRAW,
		lpfnWndProc:   C.WNDPROC(C.BridgeProc),
		hInstance:     hInstance,
		hCursor:       C.LoadCursor(hInstance, C.IDC_ARROW),
		lpszClassName: className,
	}
	fmt.Printf("%+v\n", windowClass)

	if C.RegisterClass(&windowClass) == 0 {
		fmt.Println("could not register class")
	}

	wSize := C.RECT{
		left:   0,
		top:    0,
		right:  1280,
		bottom: 720,
	}
	C.AdjustWindowRect(&wSize, C.WS_OVERLAPPEDWINDOW, C.FALSE)

	window := C.CreateWindowEx(0, className, C.CString("platform"), C.WS_OVERLAPPEDWINDOW|C.WS_VISIBLE,
		C.CW_USEDEFAULT, C.CW_USEDEFAULT, C.int(wSize.right-wSize.left), C.int(wSize.bottom-wSize.top), nil, nil, hInstance, nil)
	if window == nil {
		fmt.Println("could not create window")
	}

	resizeDIBSection(&globalBackbuffer, 1280, 720)

	p.SetUp()

	globalIsRunning = true
	for globalIsRunning {
		var message C.MSG
		for C.PeekMessage(&message, window, 0, 0, C.PM_REMOVE) != 0 {
			//fmt.Printf("%+v\n", message)

			switch message.message {
			case C.WM_QUIT:
				globalIsRunning = false

			case C.WM_SYSKEYDOWN, C.WM_SYSKEYUP, C.WM_KEYDOWN, C.WM_KEYUP:
				if message.wParam == C.VK_ESCAPE {
					globalIsRunning = false
				}

			default:
				C.TranslateMessage(&message)
				C.DispatchMessage(&message)
			}
		}

		sliceHdr := reflect.SliceHeader{
			Data: uintptr(globalBackbuffer.memory),
			Len:  globalBackbuffer.bps * (globalBackbuffer.width * globalBackbuffer.height),
		}
		sliceHdr.Cap = sliceHdr.Len

		backbuffer := BitmapBuffer{
			Memory: *(*[]uint32)(unsafe.Pointer(&sliceHdr)),
			Width:  globalBackbuffer.width,
			Height: globalBackbuffer.height,
			Bps:    globalBackbuffer.bps,
			Pitch:  globalBackbuffer.pitch,
		}
		p.UpdateAndRender(&backbuffer)

		var clientRect C.RECT
		C.GetClientRect(window, &clientRect)
		dc := C.GetDC(window)
		blitBufferInWindow(&globalBackbuffer, dc, int(clientRect.right-clientRect.left), int(clientRect.bottom-clientRect.top))
	}

	p.TearDown()

	C.DestroyWindow(window)
}