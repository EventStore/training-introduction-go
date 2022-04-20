package commands

import "time"

type Schedule struct {
	Id        string
	StartTime time.Time
	Duration  time.Duration
}

func NewScheduleCommand(id string, startTime time.Time, duration time.Duration) Schedule {
	return Schedule{
		Id:        id,
		StartTime: startTime,
		Duration:  duration,
	}
}
