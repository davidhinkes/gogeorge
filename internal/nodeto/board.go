package nodeto

import (
	"log"
	"reflect"
)

type Board struct {
	bindings []*binding
	context  context
}

func (b *Board) Once() {
	didSomething := true
	for didSomething {
		didSomething = false
		for _, binding := range b.bindings {
			if binding.done {
				continue
			}
			if !binding.canRun() {
				continue
			}
			binding.do(b.context)
			didSomething = true
		}
	}
	for _, binding := range b.bindings {
		if !binding.done {
			log.Fatalf("binding %v not complete", binding)
		}
		binding.reset()
	}
	b.context.iteration++
}

func (b Board) Run() {
	for {
		b.Once()
	}
}

func (board *Board) Bind(module Module, ins ...*Pin) []*Pin {
	inputs, outputs, err := getModuleTypes(module)
	if err != nil {
		log.Fatal(err)
	}
	if len(inputs) != len(ins) {
		log.Fatalf("number of pin inputs should match number of arguments to module: %v vs %v", len(inputs), len(ins))
	}
	n := len(inputs)
	for i := 0; i < n; i++ {
		fun := inputs[i]
		pin := ins[i].valueType
		if !pin.AssignableTo(fun) {
			log.Fatalf("module %v input %v: type %v is not assignable to %v", reflect.TypeOf(module), i, fun, pin)
		}
	}
	var outputPins []*Pin
	for _, output := range outputs {
		pin := &Pin{
			valueType: output,
		}
		outputPins = append(outputPins, pin)
	}
	board.bindings = append(board.bindings, &binding{
		module:  module,
		inputs:  ins,
		outputs: outputPins,
	})
	return outputPins
}
