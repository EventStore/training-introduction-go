package commands

import "time"

type Cancel struct {
	SlotId           string
	Reason           string
	CancellationTime time.Time
}

func NewCancelCommand(slotId, reason string, cancellationTime time.Time) Cancel {
	return Cancel{
		SlotId:           slotId,
		Reason:           reason,
		CancellationTime: cancellationTime,
	}
}
