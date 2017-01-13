// Binary plugin is the so entry point for XPlane.  The intended product
// of this is a dll/so.
package main

import (
	"unsafe"

	"github.com/davidhinkes/gogeorge/internal/xplane"
)

/*
extern float callback(float, float, int, void*);
*/
import "C"

func main() {
	// An emtpy main func is needed to make shared libraries.
}

var (
	state        *pluginState
)

type pluginState struct {
	callbackID xplane.Callback
}

//export XPluginStart
func XPluginStart(outName, outSig, outDesc *C.char) C.int {
	strcpy(unsafe.Pointer(outName), "gogeorge")
	strcpy(unsafe.Pointer(outSig), "gogeorge")
	strcpy(unsafe.Pointer(outDesc), "gogeorge")
	state = &pluginState{
		callbackID: xplane.RegisterCallback(C.callback),
	}
	return C.int(1)
}

//export XPluginStop
func XPluginStop() {
}

//export XPluginDisable
func XPluginDisable() {
	state.callbackID.Schedule(0)
}

//export XPluginEnable
func XPluginEnable() C.int {
	state.callbackID.Schedule(-1)
	return C.int(1)
}

//export XPluginReceiveMessage
func XPluginReceiveMessage(inFromWho, inMessage C.int, inParam unsafe.Pointer) {
}

//export callback
func callback(_ C.float, _ C.float, i C.int, _ unsafe.Pointer) C.float {
	return -1
}
