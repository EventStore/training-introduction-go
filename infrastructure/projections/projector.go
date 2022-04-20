package projections

import (
	"reflect"
)

type Subscription interface {
	Project(t reflect.Type, event interface{})
}

type Projector struct {
	Subscription

	projection Projection
}

func NewProjector(p Projection) *Projector {
	return &Projector{
		projection: p,
	}
}

func (p Projector) Project(t reflect.Type, event interface{}) {
	if !p.projection.CanHandle(t) {
		return
	}

	p.projection.Handle(t, event)
}
