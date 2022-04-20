package writemodal

import "fmt"

type SlotAlreadyScheduledError struct{}

func (e SlotAlreadyScheduledError) Error() string {
	return fmt.Sprintf("slot already scheduled error")
}

type SlotNotScheduledError struct{}

func (e SlotNotScheduledError) Error() string {
	return fmt.Sprintf("slot not scheduled error")
}

type SlotAlreadyBookedError struct{}

func (e SlotAlreadyBookedError) Error() string {
	return fmt.Sprintf("slot already booked error")
}

type SlotNotBookedError struct{}

func (e SlotNotBookedError) Error() string {
	return fmt.Sprintf("slot not booked error")
}

type SlotAlreadyStartedError struct{}

func (e SlotAlreadyStartedError) Error() string {
	return fmt.Sprintf("slot already started error")
}
