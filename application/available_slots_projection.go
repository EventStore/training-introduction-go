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
		s := e.(events.Scheduled)
		r.Add(readmodel.NewAvailableSlot(s.SlotId, s.StartTime, s.Duration))
	})

	p.When(events.Booked{}, func(e interface{}) {
		b := e.(events.Booked)
		r.MarkAsUnavailable(b.SlotId)
	})

	p.When(events.Cancelled{}, func(e interface{}) {
		c := e.(events.Cancelled)
		r.MarkAsAvailable(c.SlotId)
	})

	return &AvailableSlotsProjection{p, r}
}
