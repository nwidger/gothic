// +build notk

package gothic

/*
#include "interpreter.h"
*/
import "C"
import (
	"errors"
	"image"
)

func tkMainLoop() {
	for {
		C.Tcl_DoOneEvent(0)
	}
}

func tkInit(C *C.Tcl_Interp) error {
	return nil
}

func (ir *interpreter) upload_image(name string, img image.Image) error {
	return errors.New("not supported")
}
