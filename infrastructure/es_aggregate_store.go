package infrastructure

import (
	"reflect"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/EventStore/training-introduction-go/eventsourcing"
)

type EsAggregateStore struct {
	AggregateStore

	store EventStore
}

func NewEsAggregateStore(store EventStore) *EsAggregateStore {
	return &EsAggregateStore{
		store: store,
	}
}

func (s *EsAggregateStore) Save(a eventsourcing.AggregateRoot) error {
	changes := a.GetChanges()
	err := s.store.AppendEvents(s.getStreamName(a, a.GetId()), a.GetVersion(), changes...)
	if err != nil {
		return err
	}

	a.ClearChanges()
	return nil
}

func (s *EsAggregateStore) Load(aggregateId string, aggregate eventsourcing.AggregateRoot) error {
	events, err := s.store.LoadEvents(s.getStreamName(aggregate, aggregateId))
	if err != nil {
		if esdbError, ok := esdb.FromError(err); !ok {
			if esdbError.Code() == esdb.ErrorResourceNotFound {
				return &AggregateNotFoundError{}
			}
		}
		return err
	}

	aggregate.Load(events)
	aggregate.ClearChanges()
	return nil
}

func (s *EsAggregateStore) getStreamName(aggregate eventsourcing.AggregateRoot, aggregateId string) string {
	return reflect.ValueOf(aggregate).Type().String() + "-" + aggregateId
}
