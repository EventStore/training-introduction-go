package writemodal

import (
	"time"

	"github.com/EventStore/training-introduction-go/domain/writemodal/events"
	"github.com/EventStore/training-introduction-go/eventsourcing"
)

type SlotAggregate struct {
	eventsourcing.AggregateBase

	isBooked    bool
	isScheduled bool
	startTime   time.Time
}

func NewSlotAggregate() *SlotAggregate {
	a := &SlotAggregate{
		AggregateBase: eventsourcing.NewAggregateRoot(),
	}

	a.Register(events.Scheduled{}, func(e interface{}) { a.Scheduled(e.(events.Scheduled)) })
	a.Register(events.Booked{}, func(e interface{}) { a.Booked(e.(events.Booked)) })
	a.Register(events.Cancelled{}, func(e interface{}) { a.Cancelled(e.(events.Cancelled)) })
	return a
}

func (s *SlotAggregate) Schedule(id string, startTime time.Time, duration time.Duration) error {
	if s.isScheduled {
		return &SlotAlreadyScheduledError{}
	}

	s.Raise(events.NewScheduledEvent(id, startTime, duration))
	return nil
}

func (s *SlotAggregate) Scheduled(scheduled events.Scheduled) {
	s.isScheduled = true
	s.startTime = scheduled.StartTime
	s.Id = scheduled.SlotId
}

func (s *SlotAggregate) Book(patientId string) error {
	if !s.isScheduled {
		return &SlotNotScheduledError{}
	}

	if s.isBooked {
		return &SlotAlreadyBookedError{}
	}

	s.Raise(events.NewBookedEvent(s.GetId(), patientId))
	return nil
}

func (s *SlotAggregate) Booked(_ events.Booked) {
	s.isBooked = true
}

func (s *SlotAggregate) Cancel(reason string, cancellationTime time.Time) error {
	if !s.isBooked {
		return &SlotNotBookedError{}
	}

	if s.isStarted(cancellationTime) {
		return &SlotAlreadyStartedError{}
	}

	if s.isBooked && !s.isStarted(cancellationTime) {
		s.Raise(events.NewCancelledEvent(s.GetId(), reason))
		return nil
	}

	return nil
}

func (s *SlotAggregate) Cancelled(_ events.Cancelled) {
	s.isBooked = false
}

func (s *SlotAggregate) isStarted(cancellationTime time.Time) bool {
	return cancellationTime.After(s.startTime)
}
