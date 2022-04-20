package eventsourcing

import (
	"reflect"
)

type AggregateRoot interface {
	Load(events []interface{})
	ClearChanges()
	GetChanges() []interface{}
	GetId() string
	GetVersion() int64
}

type AggregateBase struct {
	AggregateRoot

	Id       string
	version  int64
	changes  []interface{}
	handlers map[reflect.Type]func(interface{})
}

func NewAggregateRoot() AggregateBase {
	return AggregateBase{
		version:  -1,
		changes:  make([]interface{}, 0),
		handlers: make(map[reflect.Type]func(interface{})),
	}
}

func (a *AggregateBase) Register(event interface{}, handler func(interface{})) {
	a.handlers[reflect.ValueOf(event).Type()] = handler
}

func (a *AggregateBase) Load(events []interface{}) {
	for _, event := range events {
		a.Raise(event)
		a.version++
	}
}

func (a *AggregateBase) Raise(event interface{}) {
	v := reflect.ValueOf(event)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if handler, exists := a.handlers[v.Type()]; exists {
		handler(event)
		a.changes = append(a.changes, event)
	}
}

func (a *AggregateBase) ClearChanges() {
	a.changes = []interface{}{}
}

func (a *AggregateBase) GetChanges() []interface{} {
	return a.changes
}

func (a *AggregateBase) GetId() string {
	return a.Id
}

func (a *AggregateBase) GetVersion() int64 {
	return a.version
}
