package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"reflect"

	"github.com/EventStore/EventStore-Client-Go/esdb"
)

type EventStore interface {
	AppendEvents(streamName string, version int64, events ...interface{}) error
	LoadEvents(streamName string) ([]interface{}, error)
}

type EsEventStore struct {
	EventStore

	esdb         *esdb.Client
	eventFactory *EventFactory
}

func NewEsEventStore(esdb *esdb.Client, f *EventFactory) *EsEventStore {
	return &EsEventStore{
		esdb:		  esdb,
		eventFactory: f,
	}
}

func (s *EsEventStore) AppendEvents(streamName string, version int64, events ...interface{}) error {
	if events == nil || len(events) == 0 {
		return nil
	}

	var eventData []esdb.EventData
	for _, e := range events {
		bytes, err := json.Marshal(e)
		if err != nil {
			panic(err)
		}

		eventData = append(eventData, esdb.EventData{
			ContentType: esdb.JsonContentType,
			EventType:   reflect.ValueOf(e).Type().String(),
			Data:        bytes,
		})
	}

	options := esdb.AppendToStreamOptions{}
	if version == -1 {
		options.ExpectedRevision = esdb.NoStream{}
	} else {
		options.ExpectedRevision = esdb.Revision(uint64(version))
	}

	_, err := s.esdb.AppendToStream(context.Background(), streamName, options, eventData...)
	return err
}

func (s *EsEventStore) LoadEvents(streamName string) ([]interface{}, error) {
	options := esdb.ReadStreamOptions{
		From:      esdb.Start{},
		Direction: esdb.Forwards,
	}
	stream, err := s.esdb.ReadStream(context.Background(), streamName, options, math.MaxInt64)
	if err != nil {
		return nil, err
	}

	events := make([]interface{}, 0)
	for {
		event, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, err
		}

		e, _, err := s.eventFactory.Create(event.Event.EventType, event.Event.Data)
		if err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}
