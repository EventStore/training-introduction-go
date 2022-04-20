package writemodal

import (
	"errors"

	"github.com/EventStore/training-introduction-go/domain/writemodal/commands"
	"github.com/EventStore/training-introduction-go/infrastructure"
)

type SlotHandlers struct {
	*infrastructure.CommandHandlerBase
}

func NewSlotHandlers(store infrastructure.AggregateStore) SlotHandlers {
	commandHandler := SlotHandlers{infrastructure.NewCommandHandler()}

	commandHandler.Register(commands.Schedule{}, func(s infrastructure.Command) error {
		ss := s.(commands.Schedule)
		aggregate := NewSlotAggregate()
		err := store.Load(ss.Id, aggregate)
		if errors.Is(err, &infrastructure.AggregateNotFoundError{}) {
			err = nil
		} else if err != nil {
			return err
		}

		err = aggregate.Schedule(ss.Id, ss.StartTime, ss.Duration)
		if err != nil {
			return err
		}

		return store.Save(aggregate)
	})

	commandHandler.Register(commands.Book{}, func(s infrastructure.Command) error {
		ss := s.(commands.Book)
		aggregate := NewSlotAggregate()
		err := store.Load(ss.SlotId, aggregate)
		if err != nil {
			return err
		}
		err = aggregate.Book(ss.PatientId)
		if err != nil {
			return err
		}

		return store.Save(aggregate)
	})

	commandHandler.Register(commands.Cancel{}, func(s infrastructure.Command) error {
		c := s.(commands.Cancel)
		aggregate := NewSlotAggregate()
		err := store.Load(c.SlotId, aggregate)
		if err != nil {
			return err
		}
		err = aggregate.Cancel(c.Reason, c.CancellationTime)
		if err != nil {
			return err
		}

		return store.Save(aggregate)
	})

	return commandHandler
}
