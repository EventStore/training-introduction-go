package readmodel

import "time"

type PatientSlot struct {
	ScheduledId string        `json:"scheduledId"`
	StartTime   time.Time     `json:"startTime"`
	Duration    time.Duration `json:"duration"`
	Status      string        `json:"status"`
}

func NewPatientSlot(scheduledId string, startTime time.Time, duration time.Duration) *PatientSlot {
	return &PatientSlot{
		ScheduledId: scheduledId,
		StartTime:   startTime,
		Duration:    duration,
		Status:      "booked",
	}
}

func (p *PatientSlot) MarkAsCancelled() {
	p.Status = "cancelled"
}
