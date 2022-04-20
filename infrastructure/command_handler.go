package infrastructure

import "reflect"

type Command interface{}

type CommandHandler interface {
	GetHandlers() map[reflect.Type]func(command Command) error
}

type CommandHandle interface {
	Handle(command Command) error
}

type CommandHandlerBase struct {
	CommandHandler

	handlers map[reflect.Type]func(command Command) error
}

func NewCommandHandler() *CommandHandlerBase {
	commandHandler := &CommandHandlerBase{}
	commandHandler.handlers = make(map[reflect.Type]func(Command) error, 0)
	return commandHandler
}

func (c *CommandHandlerBase) GetHandlers() map[reflect.Type]func(command Command) error {
	return c.handlers
}

func (c *CommandHandlerBase) Register(command interface{}, f func(Command) error) {
	c.handlers[reflect.ValueOf(command).Type()] = f
}
