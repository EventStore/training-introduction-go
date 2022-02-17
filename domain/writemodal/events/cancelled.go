package events

type Cancelled struct {
	SlotId string `json:"slotId"`
	Reason string `json:"reason"`
}

func NewCancelledEvent(slotId, reason string) Cancelled {
	return Cancelled{
		SlotId: slotId,
		Reason: reason,
	}
}
