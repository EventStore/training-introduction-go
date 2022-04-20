package projections

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ProjectionTests struct {
	t *testing.T

	projectionFactory func() Projection
}

func NewProjectionTests(t *testing.T, projectionFactory func() Projection) ProjectionTests {
	return ProjectionTests{t: t, projectionFactory: projectionFactory}
}

func (p *ProjectionTests) Given(events ...interface{}) {
	projection := p.projectionFactory()

	for _, e := range events {
		projection.Handle(reflect.ValueOf(e).Type(), e)
	}
}

func (p *ProjectionTests) Then(expected, actual interface{}) {
	assert.Equal(p.t, expected, actual)
}
