// Copyright (c) 2023, Corelight, Inc. All rights reserved.

package encoding

import (
	"fmt"
	"strings"
	"time"
)

const EventMetaDataTypeTimestamp = 1

// Event is a more convenient representation of a Zeek event
// (as opposed to an encoding.DataMessage with special contents).
type Event struct {
	Name      string
	Arguments []Data
	Metadata  []EventMetaEntry
}

// EventMetaEntry os a (id, value) tuple used to encode Zeek event metadata.
type EventMetaEntry struct {
	ID    uint64
	Value Data
}

// Encode encodes an event metadata entry as a Zeek vector.
func (e EventMetaEntry) Encode() Data {
	return Vector(
		Count(e.ID),
		e.Value,
	)
}

// String renders an event metadata entry as a string.
func (e EventMetaEntry) String() string {
	return fmt.Sprintf("%d:{%s}", e.ID, e.Value.String())
}

// NewEvent is a helper to create an encoding.Event using variadic arguments.
func NewEvent(name string, arguments ...Data) Event {
	return Event{
		Name:      name,
		Arguments: arguments,
	}
}

// Encode encodes an Event into an encoding.DataMessage given the provided topic.
func (e Event) Encode(topic string) DataMessage {
	var data Data
	if len(e.Metadata) > 0 {
		meta := make([]Data, len(e.Metadata))
		for i, m := range e.Metadata {
			meta[i] = m.Encode()
		}
		data = Vector(
			Count(1),
			Count(1),
			Vector(
				String(e.Name),
				Vector(e.Arguments...),
			),
			Vector(meta...),
		)
	} else {
		data = Vector(
			Count(1),
			Count(1),
			Vector(
				String(e.Name),
				Vector(e.Arguments...),
			),
		)
	}

	return DataMessage{
		ConstType: "data-message",
		Topic:     topic,
		Data:      &data,
	}
}

// SetTimestamp adds or replaces the event metadata timestamp (using the current time,
// if the timestamp argument is nil).
func (e *Event) SetTimestamp(timestamp *time.Time) {
	if timestamp == nil {
		now := time.Now()
		timestamp = &now
	}

	e.SetMetadata(EventMetaDataTypeTimestamp, Timestamp(*timestamp), true)
}

// SetMetadata adds or replaces the event metadata. If replace is true then
// value will be assigned to all existing entries with a matching id.
func (e *Event) SetMetadata(id uint8, value Data, replace bool) {
	em := EventMetaEntry{
		ID:    uint64(id),
		Value: value,
	}

	if !replace {
		e.Metadata = append(e.Metadata, em)
		return
	}

	for i, m := range e.Metadata {
		if m.ID == uint64(id) {
			e.Metadata[i] = em
		}
	}
}

// DeleteMetadata deletes event metadata matching id.
func (e *Event) DeleteMetadata(id uint8) {
	newMetadata := make([]EventMetaEntry, 0)

	for _, m := range e.Metadata {
		if m.ID != uint64(id) {
			newMetadata = append(newMetadata, m)
		}
	}

	e.Metadata = newMetadata
}

// String implements the Stringer interface for Event and produces a compact string representation of a zeek event.
func (e Event) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("event %s(", e.Name))
	for i, arg := range e.Arguments {
		sb.WriteString(arg.String())
		if i == len(e.Arguments)-1 {
			sb.WriteString(")")
		} else {
			sb.WriteString(", ")
		}
	}
	if len(e.Metadata) > 0 {
		sb.WriteString("[")
		for i, m := range e.Metadata {
			sb.WriteString(m.String())
			if i == len(e.Metadata)-1 {
				sb.WriteString("]")
			} else {
				sb.WriteString(", ")
			}
		}
	}
	return sb.String()
}
