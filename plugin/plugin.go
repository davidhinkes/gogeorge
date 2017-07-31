// Binary plugin is the so entry point for XPlane.  The intended product
// of this is a dll/so.
package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
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
	go http.ListenAndServe(":8080", nil)
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
	state.writeActuatorCommand(state.controller.Do(state.mkSensorData()))
	return deltaT
}

func (p *pluginState) mkSensorData() controller.SensorData {
	var roll float32
	var heading float32
	p.getDataRef("sim/cockpit2/gauges/indicators/roll_electric_deg_pilot").Get(&roll)
	p.getDataRef("sim/cockpit2/gauges/indicators/heading_electric_deg_mag_pilot").Get(&heading)
	return controller.SensorData{
		RollDegrees:    float64(roll),
		HeadingDegrees: float64(heading),
	}
}

func (p *pluginState) writeActuatorCommand(c controller.ActuatorCommand) {
	override := p.getDataRef("sim/operation/override/override_joystick_roll")
	if !c.Enabled {
		override.Set(int(0))
		return
	}
	override.Set(int(1))
	p.getDataRef("sim/joystick/yoke_roll_ratio").Set(float32(c.YokeRoll))
}
