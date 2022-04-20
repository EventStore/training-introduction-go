package infrastructure

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type EventFactory struct {
	typesByName map[string]reflect.Type
}

func NewEventFactory(eventTypes ...interface{}) *EventFactory {
	f := &EventFactory{}
	f.registerEventTypes(eventTypes)
	return f
}

func (f *EventFactory) registerEventTypes(eventTypes []interface{}) {
	f.typesByName = map[string]reflect.Type{}
	for _, e := range eventTypes {
		t := reflect.ValueOf(e).Type()
		f.typesByName[t.String()] = t
	}
}

func (f EventFactory) Create(eventType string, data []byte) (interface{}, reflect.Type, error) {
	if t, exists := f.typesByName[eventType]; exists {
		ev := reflect.New(t).Interface()
		err := json.Unmarshal(data, &ev)
		return reflect.ValueOf(ev).Elem().Interface(), t, err
	}
	return nil, nil, fmt.Errorf("unknown event type: %s", eventType)
}
