// +build notk

package gothic

/*
#include "interpreter.h"
*/
import "C"
import (
	"errors"
	"image"
	"unsafe"
)

func tkMainLoop() {

}

func tkInit(C *C.Tcl_Interp) error {
	return nil
}

func (ir *interpreter) upload_image(name string, img image.Image) error {
	return errors.New("not supported")
}

//export _gotk_go_async_handler
func _gotk_go_async_handler(ev unsafe.Pointer, flags C.int) C.int {
	return 0
}
