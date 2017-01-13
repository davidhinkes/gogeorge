package main 

import (
	//#include <stdlib.h>
	//#include <strings.h>
	"C"
	"unsafe"
)

// Strcpy is just like C's strcpy.
// Note that we deal with unsafe.Pointer instead of
// *C.char in the API because *C.char can't be exported.
func strcpy(dest unsafe.Pointer, msg string) {
	ptr := C.CString(msg)
	defer C.free(unsafe.Pointer(ptr))
	C.memcpy(dest, unsafe.Pointer(ptr), C.size_t(len(msg))+1)
}
