package readmodel

import "time"

type AvailableSlot struct {
	ScheduleId         string        `json:"scheduleId"`
	ScheduledStartTime time.Time     `json:"scheduledStartTime"`
	ScheduledDuration  time.Duration `json:"scheduledDuration"`
}

func NewAvailableSlot(id string, s time.Time, d time.Duration) AvailableSlot {
	return AvailableSlot{
		ScheduleId: id,
		ScheduledStartTime: s,
		ScheduledDuration: d,
	}
}
