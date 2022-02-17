package application

import (
	"github.com/EventStore/training-introduction-go/domain/readmodel"
	"github.com/EventStore/training-introduction-go/domain/writemodal/events"
	"github.com/EventStore/training-introduction-go/infrastructure/projections"
)

type PatientSlotsProjection struct {
	projections.ProjectionBase

	repository readmodel.PatientSlotsRepository
}

func NewPatientSlotsProjection(r readmodel.PatientSlotsRepository) *PatientSlotsProjection {
	p := projections.NewProjection()
	p.When(events.Scheduled{}, func(e interface{}) {
		s := e.(events.Scheduled)
		r.Add(readmodel.NewScheduledSlot(s.SlotId, s.StartTime, s.Duration))
	})

	p.When(events.Booked{}, func(e interface{}) {
		b := e.(events.Booked)
		r.MarkAsBooked(b.SlotId, b.PatientId)
	})

	p.When(events.Cancelled{}, func(e interface{}) {
		c := e.(events.Cancelled)
		r.MarkAsCancelled(c.SlotId)
	})

	return &PatientSlotsProjection{p, r}
}
