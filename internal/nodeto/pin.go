package nodeto

import (
	"reflect"
)

type Pin struct {
	value     interface{}
	valueType reflect.Type
	written   bool
}

func (p *Pin) reset() {
	p.value = nil
	p.written = false
}
