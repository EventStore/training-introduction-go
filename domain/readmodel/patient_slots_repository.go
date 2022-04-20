package readmodel

type PatientSlotsRepository interface {
	Add(slot ScheduledSlot)
	MarkAsBooked(scheduledSlotId, patientId string)
	MarkAsCancelled(scheduledId string)
	GetPatientSlots(patientId string) []*PatientSlot
	Clear()
}
