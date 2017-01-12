package nodeto

import (
	"log"
	"reflect"
)

var (
	bindings []binding
)

type binding struct {
	module  Module
	inputs  []*Pin
	outputs []*Pin
	done    bool
}

func (b *binding) reset() {
	b.done = false
	for _, p := range b.inputs {
		p.reset()
	}
	for _, p := range b.outputs {
		p.reset()
	}
}

func (b binding) canRun() bool {
	for _, input := range b.inputs {
		if !input.written {
			return false
		}
	}
	return true
}

func (b *binding) do(ctx Context) {
	methodValue := reflect.ValueOf(b.module).MethodByName(moduleMethodName)
	arguments := make([]reflect.Value, 0, 1+len(b.inputs))
	arguments = append(arguments, reflect.ValueOf(ctx))
	for _, pin := range b.inputs {
		arguments = append(arguments, reflect.ValueOf(pin.value))
	}
	outputValues := methodValue.Call(arguments)
	if len(outputValues) != len(b.outputs) {
		log.Fatalf("expecting number of function outputs to match pin outputs")
	}
	for i, val := range outputValues {
		b.outputs[i].value = val.Interface()
		b.outputs[i].written = true
	}
	b.done = true
}
