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

func DebugLog(msg string) {
  ptr := C.CString(msg)
	defer C.free(unsafe.Pointer(ptr))
  C.XPLMDebugString(ptr)
}
