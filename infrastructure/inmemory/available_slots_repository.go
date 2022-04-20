package inmemory

import (
	"time"

	"github.com/EventStore/training-introduction-go/domain/readmodel"
)

type InMemoryAvailableSlotsRepository struct {
	readmodel.AvailableSlotsRepository

	available []readmodel.AvailableSlot
	booked    []readmodel.AvailableSlot
}

func NewInMemoryAvailableSlotsRepository() *InMemoryAvailableSlotsRepository {
	r := &InMemoryAvailableSlotsRepository{}
	r.Clear()
	return r
}

func (r *InMemoryAvailableSlotsRepository) Add(slot readmodel.AvailableSlot) {
	r.available = append(r.available, slot)
}

func (r *InMemoryAvailableSlotsRepository) MarkAsUnavailable(availableSlotId string) {
	slots, removed := removeAvailableSlot(r.available, availableSlotId)
	r.available = slots
	r.booked = append(r.booked, removed)
}

func (r *InMemoryAvailableSlotsRepository) MarkAsAvailable(availableSlotId string) {
	slots, removed := removeAvailableSlot(r.booked, availableSlotId)
	r.booked = slots
	r.available = append(r.available, removed)
}

func removeAvailableSlot(slots []readmodel.AvailableSlot, availableSlotId string) ([]readmodel.AvailableSlot, readmodel.AvailableSlot) {
	update := make([]readmodel.AvailableSlot, 0)
	removedSlot := readmodel.AvailableSlot{}
	for _, availableSlot := range slots {
		if availableSlot.ScheduleId == availableSlotId {
			removedSlot = availableSlot
			continue
		}

		update = append(update, availableSlot)
	}

	return update, removedSlot
}

func (r *InMemoryAvailableSlotsRepository) GetSlotsAvailableOn(time time.Time) []readmodel.AvailableSlot {
	availabilityYear := time.Year()
	availabilityYearDay := time.YearDay()
	availableOnDate := make([]readmodel.AvailableSlot, 0)
	for _, slot := range r.available {
		scheduledStart := slot.ScheduledStartTime
		if availabilityYear == scheduledStart.Year() && availabilityYearDay == scheduledStart.YearDay() {
			availableOnDate = append(availableOnDate, slot)
		}
	}

	return availableOnDate
}

func (r *InMemoryAvailableSlotsRepository) Clear() {
	r.available = []readmodel.AvailableSlot{}
	r.booked = []readmodel.AvailableSlot{}
}
