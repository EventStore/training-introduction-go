package commands

type Book struct {
	SlotId    string
	PatientId string
}

func NewBookCommand(slotId, patientId string) Book {
	return Book{
		SlotId:    slotId,
		PatientId: patientId,
	}
}
