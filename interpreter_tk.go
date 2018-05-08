// +build !notk

package gothic

/*
#cgo !tcl85 LDFLAGS: -ltk8.6
#cgo tcl85 LDFLAGS: -ltk8.5
#cgo darwin tcl85 CFLAGS: -I/opt/X11/include

#include <tk.h>
#include "interpreter.h"
*/
import "C"
import (
	"bytes"
	"errors"
	"image"
	"unsafe"
)

func tkMainLoop() {
	C.Tk_MainLoop()
}

func tkInit(C *C.Tcl_Interp) error {
	status := C.Tk_Init(C)
	if status != C.TCL_OK {
		return errors.New(C.GoString(C.Tcl_GetStringResult(C)))
	}
	return nil
}

func (ir *interpreter) upload_image(name string, img image.Image) error {
	var buf bytes.Buffer
	err := sprintf(&buf, "image create photo %{}", name)
	if err != nil {
		return err
	}

	nrgba, ok := img.(*image.NRGBA)
	if !ok {
		// let's do it slowpoke
		bounds := img.Bounds()
		nrgba = image.NewNRGBA(bounds)
		for x := 0; x < bounds.Max.X; x++ {
			for y := 0; y < bounds.Max.Y; y++ {
				nrgba.Set(x, y, img.At(x, y))
			}
		}
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	handle := C.Tk_FindPhoto(ir.C, cname)
	if handle == nil {
		err := ir.eval(buf.Bytes())
		if err != nil {
			return err
		}
		handle = C.Tk_FindPhoto(ir.C, cname)
		if handle == nil {
			return errors.New("failed to create an image handle")
		}
	}

	imgdata := C.CBytes(nrgba.Pix)
	defer C.free(imgdata)

	block := C.Tk_PhotoImageBlock{
		(*C.uchar)(imgdata),
		C.int(nrgba.Rect.Max.X),
		C.int(nrgba.Rect.Max.Y),
		C.int(nrgba.Stride),
		4,
		[...]C.int{0, 1, 2, 3},
	}

	status := C.Tk_PhotoPutBlock(ir.C, handle, &block, 0, 0,
		C.int(nrgba.Rect.Max.X), C.int(nrgba.Rect.Max.Y),
		C.TK_PHOTO_COMPOSITE_SET)
	if status != C.TCL_OK {
		return errors.New(C.GoString(C.Tcl_GetStringResult(ir.C)))
	}
	return nil
}

//export _gotk_go_async_handler
func _gotk_go_async_handler(ev unsafe.Pointer, flags C.int) C.int {
	if flags != C.TK_ALL_EVENTS {
		return 0
	}
	event := (*C.GoTkAsyncEvent)(ev)
	ir := global_handles.get(int(event.go_interp)).(*interpreter)
	action := <-ir.queue
	if action.result == nil {
		action.action()
	} else {
		*action.result = action.action()
	}
	action.cond.L.Lock()
	action.cond.Signal()
	action.cond.L.Unlock()
	return 1
}
