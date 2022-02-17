package inmemory

import (
	"github.com/EventStore/training-introduction-go/domain/readmodel"
)

type InMemoryPatientSlotsRepository struct {
	readmodel.PatientSlotsRepository

	scheduledSlots []readmodel.ScheduledSlot
	patientSlots   map[string][]*readmodel.PatientSlot
}

func NewInMemoryPatientSlotsRepository() *InMemoryPatientSlotsRepository {
	r := &InMemoryPatientSlotsRepository{}
	r.Clear()
	return r
}

func (r *InMemoryPatientSlotsRepository) Add(slot readmodel.ScheduledSlot) {
	r.scheduledSlots = append(r.scheduledSlots, slot)
}

func (r *InMemoryPatientSlotsRepository) MarkAsBooked(scheduledSlotId, patientId string) {
	slots, removed := removeScheduledSlot(r.scheduledSlots, scheduledSlotId)
	patientSlot := readmodel.NewPatientSlot(removed.ScheduledSlotId, removed.StartTime, removed.Duration)
	r.scheduledSlots = slots

	if slots, exists := r.patientSlots[patientId]; exists {
		slots = append(slots, patientSlot)
	} else {
		r.patientSlots[patientId] = []*readmodel.PatientSlot{patientSlot}
	}
}

func (r *InMemoryPatientSlotsRepository) MarkAsCancelled(scheduledId string) {
	for _, patientSlots := range r.patientSlots {
		for _, patientSlot := range patientSlots {
			if patientSlot.ScheduledId == scheduledId {
				patientSlot.MarkAsCancelled()

				scheduledSlot := readmodel.NewScheduledSlot(patientSlot.ScheduledId, patientSlot.StartTime, patientSlot.Duration)
				r.scheduledSlots = append(r.scheduledSlots, scheduledSlot)
				return
			}
		}
	}
}

func removeScheduledSlot(slots []readmodel.ScheduledSlot, availableSlotId string) ([]readmodel.ScheduledSlot, readmodel.ScheduledSlot) {
	update := make([]readmodel.ScheduledSlot, 0)
	removedSlot := readmodel.ScheduledSlot{}
	for _, availableSlot := range slots {
		if availableSlot.ScheduledSlotId == availableSlotId {
			removedSlot = availableSlot
			continue
		}

		update = append(update, availableSlot)
	}

	return update, removedSlot
}

func (r *InMemoryPatientSlotsRepository) GetPatientSlots(patientId string) []*readmodel.PatientSlot {
	if slots, exists := r.patientSlots[patientId]; exists {
		return slots
	}

	return []*readmodel.PatientSlot{}
}

func (r *InMemoryPatientSlotsRepository) Clear() {
	r.scheduledSlots = []readmodel.ScheduledSlot{}
	r.patientSlots = make(map[string][]*readmodel.PatientSlot, 0)
}
