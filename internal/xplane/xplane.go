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

type Callback struct {
	id C.XPLMFlightLoopID
}

func (c Callback) Schedule(i float32) {
	C.XPLMScheduleFlightLoop(c.id, C.float(i), C.int(1))
}

func RegisterCallback(f unsafe.Pointer) Callback {
	return Callback{
		id: C.registerCallback(f),
	}
}

type DataRef struct {
	id        C.XPLMDataRef
	name      string
	writeable bool
	types     C.XPLMDataTypeID
	getter    func(C.XPLMDataRef) interface{}
	setter    func(C.XPLMDataRef, interface{})
}

func (d *DataRef) Get() interface{} {
	return d.getter(d.id)
}

func (d *DataRef) Set(value interface{}) {
	d.mustBeWritable()
	d.setter(d.id, value)
}

func getInt(id C.XPLMDataRef) interface{} {
	return int(C.XPLMGetDatai(id))
}

func getFloat32(id C.XPLMDataRef) interface{} {
	return float32(C.XPLMGetDataf(id))
}

func getFloat64(id C.XPLMDataRef) interface{} {
	return float64(C.XPLMGetDatad(id))
}

func NewDataRef(name string, zero interface{}) (*DataRef, error) {
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
	switch reflect.TypeOf(zero).Kind() {
	case reflect.Int:
		dataRef.mustSupport(C.xplmType_Int)
		dataRef.getter = getInt
    dataRef.setter = setInt
	case reflect.Float32:
		dataRef.mustSupport(C.xplmType_Float)
		dataRef.getter = getFloat32
		dataRef.setter = setFloat32
	case reflect.Float64:
		dataRef.mustSupport(C.xplmType_Double)
		dataRef.getter = getFloat64
		dataRef.setter = setFloat64
	default:
		return nil, fmt.Errorf("unspported type %v", reflect.TypeOf(zero))
	}
	return dataRef, nil
}

func (d *DataRef) mustBeWritable() {
	if d.writeable {
		return
	}
	log.Fatalf("%v is not writeable", d.name)
}

func (d *DataRef) mustSupport(t C.XPLMDataTypeID) {
	if d.types&t == 0 {
		log.Fatalf("invalid type access: %v is not a %v", d.name, t)
	}
}

func setInt(id C.XPLMDataRef, x interface{}) {
	C.XPLMSetDatai(id, C.int(x.(int)))
}

func setFloat32(id C.XPLMDataRef, x interface{}) {
	C.XPLMSetDataf(id, C.float(x.(float32)))
}

func setFloat64(id C.XPLMDataRef, x interface{}) {
	C.XPLMSetDatad(id, C.double(x.(float64)))
}

func DebugLog(msg string) {
	ptr := C.CString(msg)
	defer C.free(unsafe.Pointer(ptr))
	C.XPLMDebugString(ptr)
}
