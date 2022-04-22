package writemodal

import (
	"github.com/EventStore/training-introduction-go/domain/writemodal/events"
	"github.com/EventStore/training-introduction-go/eventsourcing"
	"time"
)

type SlotAggregate struct {
	eventsourcing.AggregateBase
}

func NewSlotAggregate() *SlotAggregate {
	a := &SlotAggregate{
		AggregateBase: eventsourcing.NewAggregateRoot(),
	}

	a.Register(events.Scheduled{}, func(e interface{}) { a.Scheduled(e.(events.Scheduled)) })
	//a.Register(events.Booked{}, func(e interface{}) { a.Booked(e.(events.Booked)) })
	//a.Register(events.Cancelled{}, func(e interface{}) { a.Cancelled(e.(events.Cancelled)) })
	return a
}

func (s *SlotAggregate) Schedule(id string, startTime time.Time, duration time.Duration) error {
	return nil
}

func (s *SlotAggregate) Scheduled(scheduled events.Scheduled) {
	s.Id = scheduled.SlotId
}
