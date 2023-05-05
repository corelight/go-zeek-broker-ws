// Copyright (c) 2023, Corelight, Inc. All rights reserved.

package encoding

import (
	"fmt"
	"strings"
)

// Event is a more convenient representation of a Zeek event
// (as opposed to an encoding.DataMessage with special contents).
type Event struct {
	Name      string
	Arguments []Data
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
	data := Vector(
		Count(1),
		Count(1),
		Vector(
			String(e.Name),
			Vector(e.Arguments...),
		),
	)
	return DataMessage{
		ConstType: "data-message",
		Topic:     topic,
		Data:      &data,
	}
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
	return sb.String()
}
