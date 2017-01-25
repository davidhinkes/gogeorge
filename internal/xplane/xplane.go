// Package xplane wraps the XPlane API.
package xplane

/*
#cgo darwin CFLAGS: -I ${SRCDIR}/XPSDK213/CHeaders -DAPL -DXPLM210
#cgo darwin LDFLAGS: -F ${SRCDIR}/XPSDK213/Libraries/Mac -framework XPLM

#include <XPLM/XPLMDataAccess.h>
#include <XPLM/XPLMPlugin.h>
#include <XPLM/XPLMProcessing.h>
#include <XPLM/XPLMUtilities.h>

#include <stdlib.h>

static XPLMFlightLoopID registerCallback(void *callback) {
	XPLMCreateFlightLoop_t arg = {
	  sizeof(XPLMCreateFlightLoop_t),
		xplm_FlightLoop_Phase_AfterFlightModel,
		callback,
		0};
  return XPLMCreateFlightLoop(&arg);
}
*/
import "C"

import (
	"fmt"
	"log"
	"reflect"
	"unsafe"
)

var (
	handlers = map[reflect.Kind]handler{
		reflect.Float32: float32Handler{},
		reflect.Int:     intHandler{},
	}
)

type handler interface {
	mask() C.XPLMDataTypeID
	get(C.XPLMDataRef) interface{}
	set(C.XPLMDataRef, interface{})
}

// Callback is a handler to an XPlane flight-loop function callback.
type Callback struct {
	id C.XPLMFlightLoopID
}

// Schedule makes a callback run.  Argument i is the interval
// at which the callbacks should be scheduled.  It can be negative.
// See the XPlane API docs for more details.
func (c Callback) Schedule(i float32) {
	C.XPLMScheduleFlightLoop(c.id, C.float(i), C.int(1))
}

// RegisterCallback creates a new callback given func f.
// Argument f must be available via externally via C.
func RegisterCallback(f unsafe.Pointer) Callback {
	return Callback{
		id: C.registerCallback(f),
	}
}

// DataRef is an XPlane data ref.
// The following types are supported:
// - int (C.int)
// - float32 (C.float)
type DataRef struct {
	id        C.XPLMDataRef
	name      string
	writeable bool
	types     C.XPLMDataTypeID
	getter    func(C.XPLMDataRef) interface{}
	setter    func(C.XPLMDataRef, interface{})
}

// Get reads a dataref.  Argument dest must be a pointer.
// Based on the type of dest, the correct XPlane get function is called.
// This will log.Fatal if attempting to read an dataref with a type
// that is not supported.
func (d *DataRef) Get(dest interface{}) {
	val := reflect.ValueOf(dest)
	ty := val.Type()
	if val.Kind() != reflect.Ptr {
		msg := fmt.Sprintf("%v needs to be a pointer", ty)
		die(msg)
	}
	handler, ok := handlers[ty.Elem().Kind()]
	if !ok {
		die(fmt.Sprintf("no handler for kind %v", ty.Elem().Kind()))
	}
	d.mustSupport(handler.mask())
	val.Elem().Set(reflect.ValueOf(handler.get(d.id)))
}

// Set sets a DataRef given value.  Based on the type of value, the correct
// XPlane API set function is used.  This will log.Fatal if the dataref
// does not support the type..
func (d *DataRef) Set(value interface{}) {
	d.mustBeWritable()
	val := reflect.ValueOf(value)
	ty := val.Type()
	handler, ok := handlers[ty.Kind()]
	if !ok {
		die(fmt.Sprintf("no handler for kind %v", ty.Kind()))
	}
	d.mustSupport(handler.mask())
	handler.set(d.id, value)
}

type intHandler struct{}

func (i intHandler) mask() C.XPLMDataTypeID {
	return C.xplmType_Int
}
func (i intHandler) get(id C.XPLMDataRef) interface{} {
	return int(C.XPLMGetDatai(id))
}
func (i intHandler) set(id C.XPLMDataRef, c interface{}) {
	C.XPLMSetDatai(id, C.int(c.(int)))
}

type float32Handler struct{}

func (i float32Handler) mask() C.XPLMDataTypeID {
	return C.xplmType_Float
}
func (i float32Handler) get(id C.XPLMDataRef) interface{} {
	return float32(C.XPLMGetDataf(id))
}
func (i float32Handler) set(id C.XPLMDataRef, c interface{}) {
	C.XPLMSetDataf(id, C.float(c.(float32)))
}

func NewDataRef(name string) (*DataRef, error) {
	nameCStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameCStr))
	id := C.XPLMFindDataRef(nameCStr)
	var zeroID C.XPLMDataRef
	if id == zeroID {
		return nil, fmt.Errorf("could not find data ref %v", name)
	}
	dataRef := &DataRef{
		id:        id,
		types:     C.XPLMGetDataRefTypes(id),
		writeable: C.XPLMCanWriteDataRef(id) > 0,
	}
	return dataRef, nil
}

func (d *DataRef) mustBeWritable() {
	if d.writeable {
		return
	}
	die(fmt.Sprintf("%v is not writeable", d.name))
}

func (d *DataRef) mustSupport(t C.XPLMDataTypeID) {
	if d.types&t == 0 {
		die(fmt.Sprintf("invalid type access: %v is not a %v", d.name, t))
	}
}

func die(msg string) {
	DebugLog(msg)
	log.Fatal(msg)
}

func DebugLog(msg string) {
	ptr := C.CString(msg)
	defer C.free(unsafe.Pointer(ptr))
	C.XPLMDebugString(ptr)
}
