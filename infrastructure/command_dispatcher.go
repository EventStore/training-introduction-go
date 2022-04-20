package infrastructure

import (
	"fmt"
	"reflect"
)

type Dispatcher struct {
	commandHandlerMap CommandHandlerMap
}

func (d Dispatcher) Dispatch(command interface{}) error {
	kind := reflect.ValueOf(command).Type()
	handler, err := d.commandHandlerMap.Get(kind)
	if err != nil {
		return fmt.Errorf("no handler registered")
	}

	return handler(command)
}

func NewDispatcher(commandHandlerMap CommandHandlerMap) Dispatcher {
	return Dispatcher{
		commandHandlerMap: commandHandlerMap,
	}
}
