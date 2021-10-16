//go:build linux

package platform

import (
	"fmt"
	"reflect"
	"unsafe"
)

/*
#cgo CFLAGS:-std=c99
#cgo LDFLAGS: -ldl -lxcb -lxcb-shm -lxcb-keysyms -lxcb-errors -lX11
#include <sys/shm.h>
#include <xcb/shm.h>
#include <xcb/xcb.h>
#include <xcb/xcb_errors.h>
#include <malloc.h>
*/
import "C"

type backbufferInfo struct {
	shmID    C.int
	shmSeg   C.xcb_shm_seg_t
	pixmapID C.xcb_pixmap_t

	memory unsafe.Pointer
	width  int32
	height int32
	bps    int32
	pitch  int32
}

func resizeBackbuffer(backbuffer *backbufferInfo, connection *C.xcb_connection_t, window C.xcb_drawable_t, width, height int32) {
	if backbuffer.shmID != 0 {
		C.xcb_shm_detach(connection, backbuffer.shmSeg)
		C.shmctl(backbuffer.shmID, C.IPC_RMID, nil)
		C.shmdt(backbuffer.memory)
		C.xcb_free_pixmap(connection, backbuffer.pixmapID)
	}
	C.xcb_flush(connection)

	backbuffer.width = width
	backbuffer.height = height
	backbuffer.bps = 4
	backbuffer.pitch = backbuffer.width * backbuffer.bps

	backbuffer.shmID = C.shmget(C.IPC_PRIVATE, C.ulong(backbuffer.bps*(width*height)), C.IPC_CREAT|0600)
	backbuffer.shmSeg = C.xcb_generate_id(connection)
	backbuffer.memory = C.shmat(backbuffer.shmID, nil, 0)

	C.xcb_shm_attach(connection, backbuffer.shmSeg, C.uint(backbuffer.shmID), 0)
	screen := C.xcb_setup_roots_iterator(C.xcb_get_setup(connection)).data
	C.xcb_shm_create_pixmap(connection, backbuffer.pixmapID, window, C.ushort(backbuffer.width), C.ushort(backbuffer.height), screen.root_depth, backbuffer.shmSeg, 0)
	C.xcb_flush(connection)
}

func (p *Platform) Main() {
	connection := C.xcb_connect(nil, nil)
	defer C.xcb_disconnect(connection)

	screen := C.xcb_setup_roots_iterator(C.xcb_get_setup(connection)).data

	var mask uint32
	var values [2]uint32

	window := C.xcb_generate_id(connection)
	mask = C.XCB_CW_EVENT_MASK
	values[0] = C.XCB_EVENT_MASK_EXPOSURE | C.XCB_EVENT_MASK_KEY_PRESS | C.XCB_EVENT_MASK_KEY_RELEASE

	C.xcb_create_window(connection, C.XCB_COPY_FROM_PARENT, window, screen.root, 0, 0, 1280, 720, 10, C.XCB_WINDOW_CLASS_INPUT_OUTPUT, screen.root_visual, C.uint(mask), unsafe.Pointer(&values[0]))
	defer C.xcb_destroy_window(connection, window)

	graphicsContext := C.xcb_generate_id(connection)
	C.xcb_create_gc(connection, graphicsContext, window, 0, nil)

	// NOTE: For getting client messages, specifically the delete window event from decorators.
	// NOTE: Also, should free string after usage.
	protocolReply := C.xcb_intern_atom_reply(connection, C.xcb_intern_atom(connection, 1, 12, C.CString("WM_PROTOCOLS")), nil)
	deleteWindowReply := C.xcb_intern_atom_reply(connection, C.xcb_intern_atom(connection, 0, 16, C.CString("WM_DELETE_WINDOW")), nil)
	C.xcb_change_property(connection, C.XCB_PROP_MODE_REPLACE, window, protocolReply.atom, C.XCB_ATOM_ATOM, 32, 1, unsafe.Pointer(&deleteWindowReply.atom))

	C.xcb_map_window(connection, window)
	C.xcb_flush(connection)

	// NOTE: SHM Support.
	shmReply := C.xcb_shm_query_version_reply(connection, C.xcb_shm_query_version(connection), nil)
	if shmReply == nil || shmReply.shared_pixmaps == 0 {
		fmt.Println("Could either not query for SHM support or lacking SHM support..")
		return
	}

	var backbuffers [2]backbufferInfo
	var currentBackbuffer *backbufferInfo

	backbuffers[0].pixmapID = C.xcb_generate_id(connection)
	backbuffers[1].pixmapID = C.xcb_generate_id(connection)
	resizeBackbuffer(&backbuffers[0], connection, window, 1280, 720)
	resizeBackbuffer(&backbuffers[1], connection, window, 1280, 720)

	shmCompletionEvent := int16(C.xcb_get_extension_data(connection, &C.xcb_shm_id).first_event + C.XCB_SHM_COMPLETION)

	p.App.SetUp()

	isRunning := true
	readyToBlit := true
	currBuffer := 0
	currentBackbuffer = &backbuffers[0]
	for isRunning {
		event := C.xcb_poll_for_event(connection)
		for ; event != nil; event = C.xcb_poll_for_event(connection) {
			//fmt.Printf("%+v\n", event)

			switch int16(event.response_type) &^ 0x80 {
			case 0:
				err := (*C.xcb_generic_error_t)(unsafe.Pointer(event))
				var errCtx *C.xcb_errors_context_t
				C.xcb_errors_context_new(connection, &errCtx)
				major := C.GoString(C.xcb_errors_get_name_for_major_code(errCtx, err.major_code))
				minor := C.GoString(C.xcb_errors_get_name_for_minor_code(errCtx, err.major_code, err.minor_code))

				var cExtension *C.char
				error := C.GoString(C.xcb_errors_get_name_for_error(errCtx, err.error_code, &cExtension))
				extension := C.GoString(cExtension)

				if minor == "" {
					minor = "no_minor"
				}
				if extension == "" {
					extension = "no_extension"
				}
				fmt.Printf("XCB Error: %s:%s, %s:%s, resource %d sequence %d", error, extension, major, minor, err.resource_id, err.sequence)
				C.xcb_errors_context_free(errCtx)

			case C.XCB_EXPOSE:
				fmt.Println("XCB_EXPOSE")

			case C.XCB_KEY_PRESS, C.XCB_KEY_RELEASE:
				fmt.Println("XCB_KEY_PRESS/RELEASE")

			case C.XCB_CLIENT_MESSAGE:
				fmt.Println("XCB_CLIENT_MESSAGE")
				clientEvent := (*C.xcb_client_message_event_t)(unsafe.Pointer(event))
				clientAtom := *(*C.uint)(unsafe.Pointer(&clientEvent.data))
				if clientAtom == deleteWindowReply.atom {
					isRunning = false
				}

			case shmCompletionEvent:
				readyToBlit = true

			default:
				fmt.Println("unmatched event")
			}

			C.free(unsafe.Pointer(event))
		}

		sliceHdr := reflect.SliceHeader{
			Data: uintptr(currentBackbuffer.memory),
			Len:  int(currentBackbuffer.bps * (currentBackbuffer.width * currentBackbuffer.height)),
		}
		sliceHdr.Cap = sliceHdr.Len
		bitmapBuffer := BitmapBuffer{
			Memory: *(*[]uint32)(unsafe.Pointer(&sliceHdr)),
			Width:  currentBackbuffer.width,
			Height: currentBackbuffer.height,
			Bps:    currentBackbuffer.bps,
			Pitch:  currentBackbuffer.pitch,
		}
		p.App.UpdateAndRender(&bitmapBuffer)

		if readyToBlit {
			readyToBlit = false
			currBuffer = (currBuffer + 1) % 2
			renderBuffer := currentBackbuffer
			currentBackbuffer = &backbuffers[currBuffer]

			C.xcb_shm_put_image(connection, window, graphicsContext, C.ushort(renderBuffer.width), C.ushort(renderBuffer.height),
				0, 0, C.ushort(renderBuffer.width), C.ushort(renderBuffer.height), 0, 0, screen.root_depth, C.XCB_IMAGE_FORMAT_Z_PIXMAP, 1, renderBuffer.shmSeg, 0)
			C.xcb_flush(connection)
		}
	}

	p.App.TearDown()

	C.xcb_shm_detach(connection, backbuffers[0].shmSeg)
	C.shmctl(backbuffers[0].shmID, C.IPC_RMID, nil)
	C.shmdt(backbuffers[0].memory)
	C.xcb_free_pixmap(connection, backbuffers[0].pixmapID)

	C.xcb_shm_detach(connection, backbuffers[1].shmSeg)
	C.shmctl(backbuffers[1].shmID, C.IPC_RMID, nil)
	C.shmdt(backbuffers[1].memory)
	C.xcb_free_pixmap(connection, backbuffers[1].pixmapID)
}
