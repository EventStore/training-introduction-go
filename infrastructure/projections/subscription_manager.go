package projections

import (
	"context"
	"reflect"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/EventStore/training-introduction-go/infrastructure"
)

type SubscriptionManager struct {
	esdb          *esdb.Client
	subscriptions []Subscription
	streamName    string
	typesByName   map[string]reflect.Type

	eventFactory *infrastructure.EventFactory
}

func NewSubscriptionManager(esdb *esdb.Client, f *infrastructure.EventFactory, subs ...Subscription) *SubscriptionManager {
	return &SubscriptionManager{
		esdb:          esdb,
		subscriptions: subs,
		eventFactory:  f,
	}
}

func (m SubscriptionManager) Start(ctx context.Context) {
	stream, err := m.esdb.SubscribeToAll(ctx, esdb.SubscribeToAllOptions{})
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			s := stream.Recv()
			if s.EventAppeared != nil {
				rawEvent := s.EventAppeared.Event
				ev, t, err := m.eventFactory.Create(rawEvent.EventType, rawEvent.Data)
				if err != nil {
					if ev != nil {
						panic(err)
					} else {
						// ignore unknown event type
						continue
					}
				}

				for _, s := range m.subscriptions {
					s.Project(t, ev)
				}
			}

			if s.SubscriptionDropped != nil {
				panic(s.SubscriptionDropped.Error)
			}
		}
	}()
}
