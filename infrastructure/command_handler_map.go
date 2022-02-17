package infrastructure

import (
	"fmt"
	"reflect"
)

type CommandHandlerMap struct {
	handlers map[reflect.Type]func(command Command) error
}

func NewCommandHandlerMap(commandHandlers ...CommandHandler) CommandHandlerMap {
	c := CommandHandlerMap{}
	c.handlers = make(map[reflect.Type]func(command Command) error, 0)

	for _, ch := range commandHandlers {
		for k, h := range ch.GetHandlers() {
			c.handlers[k] = h
		}
	}

	return c
}

func (c *CommandHandlerMap) Get(t reflect.Type) (func(command Command) error, error) {
	if handler, exists := c.handlers[t]; exists {
		return handler, nil
	}
	return nil, fmt.Errorf("handler not found!!!")
}
