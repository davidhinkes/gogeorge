// Package modules contains modules specific to the xplane API.
package xplane

import (
	"log"
	"reflect"

	"github.com/davidhinkes/gogeorge/internal/nodeto"
	"github.com/davidhinkes/gogeorge/internal/xplane"
)

type intGetter struct {
	dataRef *xplane.DataRef
}

func (g intGetter) Do() int {
	return g.dataRef.Get().(int)
}

type float32Getter struct {
	dataRef *xplane.DataRef
}

func (g float32Getter) Do() float32 {
	return g.dataRef.Get().(float32)
}

type float64Getter struct {
	dataRef *xplane.DataRef
}

func (g float64Getter) Do() float64 {
	return g.dataRef.Get().(float64)
}

func NewGetter(name string, zero interface{}) nodeto.Module {
	dataRef, err := xplane.NewDataRef(name, zero)
	if err != nil {
		log.Fatal(err)
	}
	switch reflect.TypeOf(zero).Kind() {
	case reflect.Int:
		return intGetter{dataRef}
	case reflect.Float32:
		return float32Getter{dataRef}
	case reflect.Float64:
		return float64Getter{dataRef}
	default:
		log.Fatalf("could not find hook for type %v", reflect.TypeOf(zero))
	}
	return nil
}

type setter struct {
	dataRef *xplane.DataRef
}

func (s setter) Do(value interface{}) {
	s.dataRef.Set(value)
}

func NewSetter(name string, zero interface{}) nodeto.Module {
	dataRef, err := xplane.NewDataRef(name, zero)
	if err != nil {
		log.Fatal(err)
	}
	return setter{dataRef}
}
