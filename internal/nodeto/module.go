package nodeto

import (
	"fmt"
	"reflect"
)

const (
	moduleMethodName = "Do"
)

type Module interface{}

func getModuleTypes(module Module) ([]reflect.Type, []reflect.Type, error) {
	method, found := reflect.TypeOf(module).MethodByName(moduleMethodName)
	if !found {
		return nil, nil, fmt.Errorf("method %v not found in type %v", moduleMethodName, reflect.TypeOf(module))
	}
	var inputTypes []reflect.Type
	var outputTypes []reflect.Type
	for i := 0; i < method.Type.NumIn(); i++ {
		inputTypes = append(inputTypes, method.Type.In(i))
	}
	for i := 0; i < method.Type.NumOut(); i++ {
		outputTypes = append(outputTypes, method.Type.Out(i))
	}
	ctxType := reflect.TypeOf((*Context)(nil)).Elem()
	if len(inputTypes) < 2 || inputTypes[1] != ctxType {
		return nil, nil, fmt.Errorf("the first input type must be a %v; got %v", ctxType, inputTypes)
	}
	return inputTypes[2:], outputTypes, nil
}
