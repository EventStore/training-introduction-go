package writemodal

import (
	"testing"
	"time"

	"github.com/EventStore/training-introduction-go/infrastructure"

	"github.com/EventStore/training-introduction-go/domain/writemodal/commands"
	"github.com/EventStore/training-introduction-go/domain/writemodal/events"
	"github.com/stretchr/testify/assert"
)

func TestSlotAggregate(t *testing.T) {
	store := infrastructure.NewFakeAggregateStore()
	a := SlotTests{
		AggregateTests: infrastructure.NewAggregateTests(store),
		slotId:         "slotId",
		now:            time.Now(),
		startTime:      time.Now().Add(time.Hour),
		duration:       time.Minute * 10,
		patientId:      "patient-1234",
		cancelReason:   "No longer needed",
	}

	a.RegisterHandlers(NewSlotHandlers(store))

	// Scheduling
	t.Run("ShouldBeScheduled", a.ShouldBeScheduled)
	t.Run("ShouldNotBeDoubleScheduled", a.ShouldNotBeDoubleScheduled)
	// Booking
	t.Run("ShouldBeBooked", a.ShouldBeBooked)
	t.Run("ShouldNotBeBookedIfWasNotScheduled", a.ShouldNotBeBookedIfWasNotScheduled)
	t.Run("CantBeDoubleBooked", a.CantBeDoubleBooked)
	// Cancellation
	t.Run("CanBeCancelled", a.CanBeCancelled)
	t.Run("CancelledSlotCanBeBookedAgain", a.CancelledSlotCanBeBookedAgain)
	t.Run("CantBeCancelledAfterStartTime", a.CantBeCancelledAfterStartTime)
	t.Run("CantBeCancelledIfWasNotBooked", a.CantBeCancelledIfWasNotBooked)
}

func (t *SlotTests) ShouldBeScheduled(tt *testing.T) {
	t.Given()
	t.When(commands.NewScheduleCommand(t.slotId, t.startTime, t.duration))
	t.Then(func(changes []interface{}, err error) {
		assert.Equal(tt, events.NewScheduledEvent(t.slotId, t.startTime, t.duration), changes[0])
	})
}

func (t *SlotTests) ShouldNotBeDoubleScheduled(tt *testing.T) {
	t.Given(events.NewScheduledEvent(t.slotId, t.startTime, t.duration))
	t.When(commands.NewScheduleCommand(t.slotId, t.startTime, t.duration))
	t.Then(func(_ []interface{}, err error) {
		assert.IsType(tt, &SlotAlreadyScheduledError{}, err)
	})
}

func (t *SlotTests) ShouldBeBooked(tt *testing.T) {
	t.Given(events.NewScheduledEvent(t.slotId, t.startTime, t.duration))
	t.When(commands.NewBookCommand(t.slotId, t.patientId))
	t.Then(func(changes []interface{}, err error) {
		assert.Equal(tt, events.NewBookedEvent(t.slotId, t.patientId), changes[0])
	})
}

func (t *SlotTests) ShouldNotBeBookedIfWasNotScheduled(tt *testing.T) {
	t.Given()
	t.When(commands.NewBookCommand(t.slotId, t.patientId))
	t.Then(func(_ []interface{}, err error) {
		assert.IsType(tt, &SlotNotScheduledError{}, err)
	})
}

func (t *SlotTests) CantBeDoubleBooked(tt *testing.T) {
	t.Given(
		events.NewScheduledEvent(t.slotId, t.startTime, t.duration),
		events.NewBookedEvent(t.slotId, t.patientId))
	t.When(commands.NewBookCommand(t.slotId, t.patientId))
	t.Then(func(_ []interface{}, err error) {
		assert.IsType(tt, &SlotAlreadyBookedError{}, err)
	})
}

func (t *SlotTests) CanBeCancelled(tt *testing.T) {
	t.Given(
		events.NewScheduledEvent(t.slotId, t.startTime, t.duration),
		events.NewBookedEvent(t.slotId, t.patientId))
	t.When(commands.NewCancelCommand(t.slotId, t.cancelReason, t.now))
	t.Then(func(changes []interface{}, _ error) {
		assert.Equal(tt, events.NewCancelledEvent(t.slotId, t.cancelReason), changes[0])
	})
}

func (t *SlotTests) CancelledSlotCanBeBookedAgain(tt *testing.T) {
	t.Given(
		events.NewScheduledEvent(t.slotId, t.startTime, t.duration),
		events.NewBookedEvent(t.slotId, t.patientId),
		events.NewCancelledEvent(t.slotId, t.cancelReason))
	t.When(commands.NewBookCommand(t.slotId, t.patientId))
	t.Then(func(changes []interface{}, _ error) {
		assert.Equal(tt, events.NewBookedEvent(t.slotId, t.patientId), changes[0])
	})
}

func (t *SlotTests) CantBeCancelledAfterStartTime(tt *testing.T) {
	t.Given(
		events.NewScheduledEvent(t.slotId, t.startTime, t.duration),
		events.NewBookedEvent(t.slotId, t.patientId))
	t.When(commands.NewCancelCommand(t.slotId, t.cancelReason, t.startTime.Add(time.Hour)))
	t.Then(func(_ []interface{}, err error) {
		assert.IsType(tt, &SlotAlreadyStartedError{}, err)
	})
}

func (t *SlotTests) CantBeCancelledIfWasNotBooked(tt *testing.T) {
	t.Given(events.NewScheduledEvent(t.slotId, t.startTime, t.duration))
	t.When(commands.NewCancelCommand(t.slotId, t.cancelReason, t.now))
	t.Then(func(_ []interface{}, err error) {
		assert.IsType(tt, &SlotNotBookedError{}, err)
	})
}

type SlotTests struct {
	infrastructure.AggregateTests

	slotId       string
	patientId    string
	cancelReason string
	now          time.Time
	startTime    time.Time
	duration     time.Duration
}
