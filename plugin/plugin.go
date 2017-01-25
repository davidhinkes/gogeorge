// Binary plugin is the so entry point for XPlane.  The intended product
// of this is a dll/so.
package main

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"unsafe"

	"github.com/davidhinkes/gogeorge/controller"
	"github.com/davidhinkes/gogeorge/internal/xplane"
)

/*
extern float callback(float, float, int, void*);
*/
import "C"

const (
	deltaT = 0.1
)

func main() {
	// An emtpy main func is needed to make shared libraries.
}

var (
	state *pluginState
)

type pluginState struct {
	callbackID xplane.Callback
	roll       float64
	lock       sync.Mutex
	controller *controller.T
	dataRefs   map[string]*xplane.DataRef
}

func (p *pluginState) getDataRef(name string) *xplane.DataRef {
	if ref, ok := p.dataRefs[name]; ok {
		return ref
	}
	ref, err := xplane.NewDataRef(name)
	if err != nil {
		xplane.DebugLog(err.Error())
		log.Fatal(err)
	}
	return ref
}

//export XPluginStart
func XPluginStart(outName, outSig, outDesc *C.char) C.int {
	strcpy(unsafe.Pointer(outName), "gogeorge")
	strcpy(unsafe.Pointer(outSig), "gogeorge")
	strcpy(unsafe.Pointer(outDesc), "gogeorge")
	state = &pluginState{
		callbackID: xplane.RegisterCallback(C.callback),
		controller: controller.New(deltaT),
		dataRefs:   make(map[string]*xplane.DataRef),
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
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("recovered: %v\n stack: %s", r, debug.Stack())
			xplane.DebugLog(msg)
			log.Fatal(msg)
		}
	}()
	state.lock.Lock()
	defer state.lock.Unlock()
	state.controller.Do(state.mkSensorData())
	return deltaT
}

func (p *pluginState) mkSensorData() controller.SensorData {
	var roll float32
	p.getDataRef("sim/cockpit2/gauges/indicators/heading_electric_deg_mag_pilot").Get(&roll)
	return controller.SensorData{
		RollDegrees: float64(roll),
	}
	// p.getDataRef("sim/cockpit2/gauges/indicators/heading_electric_deg_mag_pilot")
}

func (p *pluginState) writeActuatorCommand(c controller.ActuatorCommand) {
}
