// Copyright (c) 2023, Corelight, Inc. All rights reserved.

package encoding

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// DataMessage handles the encoding of "data-message" structures (which are used to represent events and errors).
type DataMessage struct {
	ConstType string `json:"type"` // always "data-message"
	Topic     string `json:"topic"`
	Data      *Data
}

// DataMessageUnknownTypeError is raised when we receive a DataMessage that is neither an event or an error from broker.
type DataMessageUnknownTypeError struct {
	TypeValue string
}

// Error implements the Error interface for DataMessageUnknownTypeError.
func (e DataMessageUnknownTypeError) Error() string {
	return fmt.Sprintf("the DataMessage's \"type\" property is an unknown value of \"%s\"", e.TypeValue)
}

// MarshalJSON implements the Marshaler interface for DataMessage.
func (d DataMessage) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(map[string]interface{}{
		"type":       d.ConstType,
		"topic":      d.Topic,
		"@data-type": d.Data.DataType,
		"data":       d.Data.DataValue,
	}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalJSON implements the Unmarshaler interface for DataMessage
func (d *DataMessage) UnmarshalJSON(b []byte) error {
	var rawData map[string]interface{}
	if err := json.Unmarshal(b, &rawData); err != nil {
		return err
	}

	t, ok := rawData["type"]
	if !ok {
		return fmt.Errorf("DataMessage is missing the \"type\" property")
	}

	ts, ok := t.(string)
	if !ok {
		return fmt.Errorf("the DataMessage's \"type\" property is not a string")
	}

	if ts != "data-message" {
		if ts != "error" {
			return DataMessageUnknownTypeError{TypeValue: ts}
		}

		c, ok := rawData["code"]
		if !ok {
			return fmt.Errorf("ErrorMessage is missing the \"code\" property")
		}

		cs, ok := c.(string)
		if !ok {
			return fmt.Errorf("the ErrorMessage's \"code\" property is not a string")
		}

		co, ok := rawData["context"]
		if !ok {
			return fmt.Errorf("ErrorMessage is missing the \"context\" property")
		}

		cos, ok := co.(string)
		if !ok {
			return fmt.Errorf("the ErrorMessage's \"context\" property is not a string")
		}

		return ErrorMessage{
			ConstType: "error",
			Code:      cs,
			Context:   cos,
		}
	}

	d.ConstType = ts

	t, ok = rawData["topic"]
	if !ok {
		return fmt.Errorf("DataMessage is missing the \"topic\" property")
	}

	ts, ok = t.(string)
	if !ok {
		return fmt.Errorf("the DataMessage's \"topic\" property is not a string")
	}

	d.Topic = ts

	return json.Unmarshal(b, &d.Data)
}

const eventToplevelVectorLen = 3
const eventSignatureVectorLen = 2

// GetEvent obtains the topic, and Event from a zeek broker event encoded in a DataMessage.
func (d *DataMessage) GetEvent() (topic string, evt Event, err error) {
	if d.Data.DataType != TypeVector {
		return "", Event{},
			fmt.Errorf("expected data type for event to be a vector but got a %s instead",
				d.Data.DataType.String())
	}

	vec, ok := d.Data.DataValue.([]Data)
	if !ok {
		return "", Event{}, fmt.Errorf("vector value has invalid type")
	}

	if len(vec) != eventToplevelVectorLen {
		return "", Event{}, fmt.Errorf("vector value has invalid length (%d)", len(vec))
	}

	if vec[0].DataType != TypeCount {
		return "", Event{},
			fmt.Errorf("event format number has invalid type (%s)", vec[0].DataType.String())
	}

	formatNumber, ok := vec[0].DataValue.(uint64)
	if !ok {
		return "", Event{}, fmt.Errorf("event format number has invalid type")
	}

	if formatNumber != 1 {
		return "", Event{},
			fmt.Errorf("event format number has invalid value (%d", formatNumber)
	}

	if vec[1].DataType != TypeCount {
		return "", Event{},
			fmt.Errorf("event message type has invalid type (%s)", vec[1].DataType.String())
	}

	zeekMessageType, ok := vec[1].DataValue.(uint64)
	if !ok {
		return "", Event{}, fmt.Errorf("event message type has invalid type")
	}

	if zeekMessageType != 1 {
		return "", Event{},
			fmt.Errorf("event message type has invalid value (%d", zeekMessageType)
	}

	if vec[2].DataType != TypeVector {
		return "", Event{},
			fmt.Errorf("event signature has invalid type (%s)", vec[2].DataType.String())
	}

	sig, ok := vec[2].DataValue.([]Data)
	if !ok {
		return "", Event{}, fmt.Errorf("event signature has invalid type")
	}

	if len(sig) < 1 { // TODO: can events have zero arguments?
		return "", Event{}, fmt.Errorf("event signature is empty")
	}

	if sig[0].DataType != TypeString {
		return "", Event{},
			fmt.Errorf("event name has invalid type (%s)", sig[0].DataType.String())
	}

	evt.Name, ok = sig[0].DataValue.(string)
	if !ok {
		return "", Event{}, fmt.Errorf("event name has invalid type")
	}

	topic = d.Topic

	if len(sig) < eventSignatureVectorLen {
		return "", Event{}, fmt.Errorf("event has too few signature elements")
	}

	if sig[1].DataType != TypeVector {
		return "", Event{},
			fmt.Errorf("event arguments has invalid type (%s)", sig[1].DataType.String())
	}

	evt.Arguments, ok = sig[1].DataValue.([]Data)
	if !ok {
		return "", Event{}, fmt.Errorf("event arguments has invalid type")
	}

	if len(sig) > eventSignatureVectorLen {
		if sig[2].DataType != TypeVector {
			return "", Event{}, fmt.Errorf("event metadata has invalid encoded type (%s)", sig[2].DataType.String())
		}

		metadata, typeOk := sig[2].DataValue.([]Data)
		if !typeOk {
			return "", Event{}, fmt.Errorf("event metadata has invalid parsed type (%T)", metadata)
		}

		metaList := make([]EventMetaEntry, len(metadata))
		for i, metadataEntry := range metadata {
			if metadataEntry.DataType != TypeVector {
				return "", Event{}, fmt.Errorf("event metadata entry %d has invalid encoded type (%s)",
					i, metadataEntry.DataType.String())
			}

			entryVector, typeOkOk := metadataEntry.DataValue.([]Data)
			if !typeOkOk {
				return "", Event{}, fmt.Errorf("event metadata entry %d has invalid parsed type (%T)", i, metadata)
			}

			if len(entryVector) != 2 {
				return "", Event{}, fmt.Errorf("event metadata entry %d has incorrect length %d",
					i, len(entryVector))
			}

			if entryVector[0].DataType != TypeCount {
				return "", Event{}, fmt.Errorf("event metadata entry %d type ID has invalid encoded type (%s)",
					i, entryVector[0].DataType.String())
			}

			entryTypeID, typeOkOk := entryVector[0].DataValue.(uint64)
			if !typeOkOk {
				return "", Event{}, fmt.Errorf("event metadata entry %d type ID has invalid parsed type (%T)",
					i, metadata)
			}

			metaList[i].ID = entryTypeID
			metaList[i].Value = entryVector[1]
		}
		evt.Metadata = metaList
	}

	return
}
