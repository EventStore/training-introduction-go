package events

type Booked struct {
	SlotId    string `json:"slotId"`
	PatientId string `json:"patientId"`
}

func NewBookedEvent(slotId, patientId string) Booked {
	return Booked{
		SlotId: slotId,
		PatientId: patientId,
	}
}
