package io

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

var (
	variables map[string]ty
	lock      sync.Mutex
)

const (
	prefix = "/_io/"
)

func init() {
	variables = make(map[string]ty)
	http.HandleFunc(prefix, list)
}

type Bool struct {
	value bool
	lock  sync.Mutex
}

func NewBool(name string, d bool) *Bool {
	i := &Bool{}
	i.Set(d)
	add(name, i)
	return i
}

func (i *Bool) Get() bool {
	i.lock.Lock()
	defer i.lock.Unlock()
	return i.value
}

func (i *Bool) Set(val bool) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.value = val
}

func (i Bool) String() string {
	return fmt.Sprintf("%v", i.Get())
}

func (i *Bool) FromString(str string) error {
	x, err := strconv.ParseBool(str)
	if err != nil {
		return err
	}
	i.Set(bool(x))
	return nil
}

type Int struct {
	value int
	lock  sync.Mutex
}

func NewInt(name string, d int) *Int {
	i := &Int{}
	i.Set(d)
	add(name, i)
	return i
}

func (i *Int) Get() int {
	i.lock.Lock()
	defer i.lock.Unlock()
	return i.value
}

func (i *Int) Set(val int) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.value = val
}

func (i Int) String() string {
	return fmt.Sprintf("%v", i.Get())
}

func (i *Int) FromString(str string) error {
	x, err := strconv.ParseInt(str, 0, 32)
	if err != nil {
		return err
	}
	i.Set(int(x))
	return nil
}

type Float64 struct {
	value float64
	lock  sync.Mutex
}

func NewFloat64(name string, d float64) *Float64 {
	i := &Float64{}
	i.Set(d)
	add(name, i)
	return i
}

func (i *Float64) Get() float64 {
	i.lock.Lock()
	defer i.lock.Unlock()
	return i.value
}

func (i *Float64) Set(val float64) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.value = val
}

func (i *Float64) String() string {
	return fmt.Sprintf("%v", i.Get())
}

func (i *Float64) FromString(str string) error {
	x, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}
	i.Set(float64(x))
	return nil
}

type ty interface {
	String() string
	FromString(value string) error
}

type handle struct {
	ty ty
}

func (h handle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	val := req.URL.Query().Get("value")
	if val != "" {
		if err := h.ty.FromString(val); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
	buf := bytes.NewBuffer(nil)
	buf.WriteString(h.ty.String())
	w.Write(buf.Bytes())
}

func list(w http.ResponseWriter, req *http.Request) {
	lock.Lock()
	defer lock.Unlock()
	buf := bytes.NewBuffer(nil)
	for name, variable := range variables {
		buf.WriteString(fmt.Sprintf("%v\t%v\n", name, variable.String()))
	}
	w.Write(buf.Bytes())
}

func add(name string, t ty) {
	lock.Lock()
	defer lock.Unlock()
	variables[name] = t
	http.Handle(fmt.Sprintf("%v%v", prefix, name), handle{ty: t})
}
