package projections

import "reflect"

type Projection interface {
	CanHandle(reflect.Type) bool
	Handle(reflect.Type, interface{})
	GetHandledTypes() []reflect.Type
}

type ProjectionBase struct {
	Projection

	handlers     []EventHandler
	handledTypes map[reflect.Type]bool
	types        []reflect.Type
}

func NewProjection() ProjectionBase {
	return ProjectionBase{
		handledTypes: make(map[reflect.Type]bool, 0),
	}
}

func (p *ProjectionBase) When(event interface{}, handler func(interface{})) {
	t := reflect.ValueOf(event).Type()
	p.handlers = append(p.handlers, NewEventHandler(t, handler))
	p.handledTypes[t] = true
}

func (p *ProjectionBase) CanHandle(t reflect.Type) bool {
	_, exists := p.handledTypes[t]
	return exists
}

func (p *ProjectionBase) Handle(t reflect.Type, event interface{}) {
	for _, h := range p.handlers {
		if h.Type == t {
			h.Handler(event)
		}
	}
}

type EventHandler struct {
	Type    reflect.Type
	Handler func(interface{})
}

func NewEventHandler(t reflect.Type, handler func(interface{})) EventHandler {
	return EventHandler{
		Type:    t,
		Handler: handler,
	}
}
