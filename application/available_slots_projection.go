package application

import (
	"github.com/EventStore/training-introduction-go/domain/readmodel"
	"github.com/EventStore/training-introduction-go/domain/writemodal/events"
	"github.com/EventStore/training-introduction-go/infrastructure/projections"
)

type AvailableSlotsProjection struct {
	projections.ProjectionBase

	repository readmodel.AvailableSlotsRepository
}

func NewAvailableSlotsProjection(r readmodel.AvailableSlotsRepository) *AvailableSlotsProjection {
	p := projections.NewProjection()
	p.When(events.Scheduled{}, func(e interface{}) {

	})

	return &AvailableSlotsProjection{p, r}
}
