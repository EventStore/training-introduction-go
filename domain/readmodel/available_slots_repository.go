package readmodel

import "time"

type AvailableSlotsRepository interface {
	Add(AvailableSlot)
	MarkAsUnavailable(string)
	MarkAsAvailable(string)
	GetSlotsAvailableOn(time time.Time) []AvailableSlot
	Clear()
}
