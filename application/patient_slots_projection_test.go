package application

import (
	"testing"
	"time"

	"github.com/EventStore/training-introduction-go/domain/readmodel"
	"github.com/EventStore/training-introduction-go/domain/writemodal/events"
	"github.com/EventStore/training-introduction-go/infrastructure/inmemory"
	"github.com/EventStore/training-introduction-go/infrastructure/projections"
)

func TestPatientSlotsProjection(t *testing.T) {
	r := inmemory.NewInMemoryPatientSlotsRepository()
	p := &PatientSlotsTests{
		ProjectionTests: projections.NewProjectionTests(t, func() projections.Projection {
			return NewPatientSlotsProjection(r)
		}),
		repository: r,
		slotId:     "slotId",
		startTime:  time.Now(),
		duration:   time.Minute * 10,
		patientId:  "patient-1234",
		reason:     "No longer needed",
	}

	t.Run("ShouldReturnAnEmptyList", p.ShouldReturnAnEmptyList)
	t.Run("ShouldReturnAnEmptyListIfTheSlotWasScheduled", p.ShouldReturnAnEmptyListIfTheSlotWasScheduled)
	t.Run("ShouldReturnASlotIfWasBooked", p.ShouldReturnASlotIfWasBooked)
	t.Run("ShouldReturnASlotIfWasCancelled", p.ShouldReturnASlotIfWasCancelled)
	t.Run("ShouldAllowToBookPreviouslyCancelledSlot", p.ShouldAllowToBookPreviouslyCancelledSlot)
}

func (p *PatientSlotsTests) ShouldReturnAnEmptyList(t *testing.T) {
	p.Given()
	p.Then(
		[]*readmodel.PatientSlot{},
		p.repository.GetPatientSlots(p.patientId))
}

func (p *PatientSlotsTests) ShouldReturnAnEmptyListIfTheSlotWasScheduled(t *testing.T) {
	p.Given(
		events.NewScheduledEvent(p.slotId, p.startTime, p.duration))
	p.Then(
		[]*readmodel.PatientSlot{},
		p.repository.GetPatientSlots(p.patientId))
}

func (p *PatientSlotsTests) ShouldReturnASlotIfWasBooked(t *testing.T) {
	p.Given(
		events.NewScheduledEvent(p.slotId, p.startTime, p.duration),
		events.NewBookedEvent(p.slotId, p.patientId))
	p.Then(
		[]*readmodel.PatientSlot{readmodel.NewPatientSlot(p.slotId, p.startTime, p.duration)},
		p.repository.GetPatientSlots(p.patientId))
}

func (p *PatientSlotsTests) ShouldReturnASlotIfWasCancelled(t *testing.T) {
	expected := readmodel.NewPatientSlot(p.slotId, p.startTime, p.duration)
	expected.MarkAsCancelled()

	p.Given(
		events.NewScheduledEvent(p.slotId, p.startTime, p.duration),
		events.NewBookedEvent(p.slotId, p.patientId),
		events.NewCancelledEvent(p.slotId, p.reason))
	p.Then(
		[]*readmodel.PatientSlot{expected},
		p.repository.GetPatientSlots(p.patientId))
}

func (p *PatientSlotsTests) ShouldAllowToBookPreviouslyCancelledSlot(t *testing.T) {
	expected := readmodel.NewPatientSlot(p.slotId, p.startTime, p.duration)
	expected.MarkAsCancelled()

	p.Given(
		events.NewScheduledEvent(p.slotId, p.startTime, p.duration),
		events.NewBookedEvent(p.slotId, p.patientId),
		events.NewCancelledEvent(p.slotId, p.reason))
	p.Then(
		[]*readmodel.PatientSlot{expected},
		p.repository.GetPatientSlots(p.patientId))
}

type PatientSlotsTests struct {
	projections.ProjectionTests

	repository readmodel.PatientSlotsRepository
	slotId     string
	patientId  string
	reason     string
	startTime  time.Time
	duration   time.Duration
}
