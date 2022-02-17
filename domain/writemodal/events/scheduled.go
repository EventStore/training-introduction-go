package events

import (
	"time"
)

type Scheduled struct {
	SlotId    string        `json:"slotId"`
	StartTime time.Time     `json:"startTime"`
	Duration  time.Duration `json:"duration"`
}

func NewScheduledEvent(slotId string, start time.Time, duration time.Duration) Scheduled {
	return Scheduled{
		SlotId: slotId,
		StartTime: start,
		Duration: duration,
	}
}
