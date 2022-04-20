package application

import (
	"testing"
	"time"

	"github.com/EventStore/training-introduction-go/domain/readmodel"
	"github.com/EventStore/training-introduction-go/domain/writemodal/events"
	"github.com/EventStore/training-introduction-go/infrastructure/inmemory"
	"github.com/EventStore/training-introduction-go/infrastructure/projections"
)

func TestAvailableSlotsProjection(t *testing.T) {
	r := inmemory.NewInMemoryAvailableSlotsRepository()
	p := &AvailableSlotsTests{
		ProjectionTests: projections.NewProjectionTests(t, func() projections.Projection {
			return NewAvailableSlotsProjection(r)
		}),
		repository: r,
		slotId:     "slotId",
		startTime:  time.Now(),
		duration:   time.Minute * 10,
		patientId:  "patient-1234",
		reason:     "No longer needed",
	}

	t.Run("ShouldAddSlotToTheList", p.ShouldAddSlotToTheList)
	t.Run("ShouldRemoveSlotFromTheListIfWasBooked", p.ShouldRemoveSlotFromTheListIfWasBooked)
	t.Run("ShouldAddSlotAgainIfBookingWasCancelled", p.ShouldAddSlotAgainIfBookingWasCancelled)
}

func (p *AvailableSlotsTests) ShouldAddSlotToTheList(t *testing.T) {
	p.Given(
		events.NewScheduledEvent(p.slotId, p.startTime, p.duration))
	p.Then(
		[]readmodel.AvailableSlot{readmodel.NewAvailableSlot(p.slotId, p.startTime, p.duration)},
		p.repository.GetSlotsAvailableOn(p.startTime))
}

func (p *AvailableSlotsTests) ShouldRemoveSlotFromTheListIfWasBooked(t *testing.T) {
	p.Given(
		events.NewScheduledEvent(p.slotId, p.startTime, p.duration),
		events.NewBookedEvent(p.slotId, p.patientId))
	p.Then(
		[]readmodel.AvailableSlot{},
		p.repository.GetSlotsAvailableOn(p.startTime))
}

func (p *AvailableSlotsTests) ShouldAddSlotAgainIfBookingWasCancelled(t *testing.T) {

}

type AvailableSlotsTests struct {
	projections.ProjectionTests

	repository readmodel.AvailableSlotsRepository
	slotId     string
	patientId  string
	reason     string
	startTime  time.Time
	duration   time.Duration
}
