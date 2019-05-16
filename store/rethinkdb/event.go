package rethinkdb

import (
	"plateau/event"
	"plateau/store"
)

// EventContainer ...
type EventContainer struct {
	event.Event `rethinkdb:"event"`

	Emitter   Player   `rethinkdb:"emitter_id,reference" rethinkdb_ref:"id"`
	Receivers []Player `rethinkdb:"receiver_ids,reference" rethinkdb_ref:"id"`
	Subjects  []Player `rethinkdb:"subject_ids,reference" rethinkdb_ref:"id"`

	Payload map[string]interface{} `rethinkdb:"payload"`
}

func eventContainerFromStoreStruct(ec store.EventContainer) *EventContainer {
	var (
		emitter             Player
		receivers, subjects []Player
	)

	if ec.Emitter != nil {
		emitter = *playerFromStoreStruct(*ec.Emitter)
	}

	for _, p := range ec.Receivers {
		receivers = append(receivers, *playerFromStoreStruct(*p))
	}

	for _, p := range ec.Subjects {
		subjects = append(subjects, *playerFromStoreStruct(*p))
	}

	return &EventContainer{
		Event:     ec.Event,
		Emitter:   emitter,
		Receivers: receivers,
		Subjects:  subjects,
		Payload:   ec.Payload,
	}
}

func (s *EventContainer) toStoreStruct() *store.EventContainer {
	var receivers, subjects []*store.Player

	for _, p := range s.Receivers {
		receivers = append(receivers, p.toStoreStruct())
	}

	for _, p := range s.Subjects {
		subjects = append(subjects, p.toStoreStruct())
	}

	return &store.EventContainer{
		Event:     s.Event,
		Emitter:   s.Emitter.toStoreStruct(),
		Receivers: receivers,
		Subjects:  subjects,
		Payload:   s.Payload,
	}
}
