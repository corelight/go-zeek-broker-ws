// Copyright (c) 2023, Corelight, Inc. All rights reserved.

package encoding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"time"
)

//go:generate go-enum --marshal

// Type represents the types specific to Zeek broker websocket encoding
/*
ENUM(
Boolean = "boolean" // Native JSON boolean (maps to bool)
Count = "count" // 64 bit unsigned integer (maps to uint64)
Integer = "integer" // Native JSON (signed) integer (maps to int64)
Real = "real" // Native JSON "number" type (maps to float64)
Timespan = "timespan" // String-encoded time span (maps to time.Duration)
Timestamp = "timestamp" // ISO 8601 encoded time in YYYY-MM-DDThh:mm:ss.sss format (maps to time.Time)
String = "string" // Native JSON string (maps to string)
EnumValue = "enum-value" // Zeek enum value mapped to native JSON string (maps to string)
Address = "address" // String-encoded IPv4/IPv6 address (maps to net.Addr)
Subnet = "subnet" // String-encoded IPv4/IPv6 subnet in <address>/<prefix-length> format (maps to net.IPNet)
Port = "port" // String-encoded service port in <port>/<protocol> format (maps to encoding.Service)
Vector = "vector" // Sequence of encoding.Data (maps to []Data)
Set = "set" // Sequence of encoding.Data with distinct objects (maps to map[Data]struct{})
Table = "table" // Map of encoding.Data keys to encoding.Data values (maps to map[Data]Data)
None = "none" // JSON empty object, maps to nil
)
*/
type Type string

const nanosecondsIn24Hours = 8.64e+13

const brokerTimeFormat = "2006-01-02T15:04:05.000"

// Data is the recursive type/value structure used by the Zeek broker websocket encoding.
type Data struct {
	DataType  Type        `json:"@data-type"`
	DataValue interface{} `json:"data"`
}

// AckMessage is the handshake sent by broker on connect.
type AckMessage struct {
	ConstType    string `json:"type"` // always "ack"
	EndpointUUID string `json:"endpoint"`
	Version      string `json:"version"`
}

// ErrorMessage encodes error messages from broker.
type ErrorMessage struct { //nolint:errname // This is an error message first, and an Error second.
	ConstType string `json:"type"` // always "error"
	Code      string `json:"code"`
	Context   string `json:"context"`
}

// Error implements the error interface for ErrorMessage
func (e ErrorMessage) Error() string {
	return fmt.Sprintf("broker error code=\"%s\" context=\"%s\"", e.Code, e.Context)
}

// decode unpacks values from a JSON object deserialised to a map.
//
//nolint:funlen // it just needs to be long due to the verbosity of error checking
//nolint:gocognit // shush
func (d *Data) decode(rawDataPtr *map[string]interface{}) error {
	rawData := *rawDataPtr
	t, ok := rawData["@data-type"]
	if !ok {
		return fmt.Errorf("JSON object is missing the \"@data-type\" property - not a valid data object")
	}

	ts, ok := t.(string)
	if !ok {
		return fmt.Errorf("the \"@data-type\" property is not a string - not a valid data object")
	}

	v, ok := rawData["data"]
	if !ok {
		return fmt.Errorf("JSON object is missing the \"data\" property - not a valid data object")
	}

	var err error
	d.DataType, err = ParseType(ts)
	if err != nil {
		return err
	}

	var numberValue json.Number
	if d.DataType == TypeReal || d.DataType == TypeInteger || d.DataType == TypeCount {
		numberValue, ok = v.(json.Number)
		if !ok {
			return fmt.Errorf("expected %s type to be serialized as JSON number but got type %T value %v",
				d.DataType.String(), v, v)
		}
	}

	var stringValue string
	if d.DataType == TypeString || d.DataType == TypeTimespan || d.DataType == TypeTimestamp ||
		d.DataType == TypeEnumValue || d.DataType == TypeAddress || d.DataType == TypeSubnet || d.DataType == TypePort {
		stringValue, ok = v.(string)
		if !ok {
			return fmt.Errorf("expected Count type to be serialized as JSON string but got type %T value %v",
				stringValue, stringValue)
		}
	}

	switch d.DataType {
	case TypeBoolean:
		boolValue, ok := v.(bool)
		if !ok {
			return fmt.Errorf("expected Boolean type to be serialized as JSON boolean but got type %T value %v",
				boolValue, boolValue)
		}
		d.DataValue = boolValue
	case TypeNone:
		// TODO: how do we test this end-to-end with zeek?
		m, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected None type elements to be serialized empty JSON objects but got type %T value %v",
				v, v)
		}
		if len(m) > 0 {
			return fmt.Errorf("expected None type to be serialized as empty JSON object it has len %d", len(m))
		}
		d.DataValue = nil
	case TypeReal:
		f, err := numberValue.Float64()
		if err != nil {
			return fmt.Errorf("problem converting Real type to float64: %w", err)
		}
		d.DataValue = f
	case TypeString:
		fallthrough
	case TypeEnumValue:
		d.DataValue = stringValue
	case TypeCount:
		i, err := strconv.ParseUint(numberValue.String(), 10, 64)
		if err != nil {
			return fmt.Errorf("problem converting Count type to uint64: %w", err)
		}
		d.DataValue = i
	case TypeInteger:
		i, err := numberValue.Int64()
		if err != nil {
			return fmt.Errorf("problem converting Integer type to int64: %w", err)
		}
		d.DataValue = i
	case TypeTimespan:
		if strings.HasSuffix(stringValue, "min") {
			stringValue = stringValue[:len(stringValue)-2]
		}
		if strings.HasSuffix(stringValue, "d") {
			days, err := strconv.ParseFloat(stringValue[:len(stringValue)-1], 64)
			if err != nil {
				return err
			}
			d.DataValue = time.Duration(math.Trunc(days * nanosecondsIn24Hours))
			return nil
		}
		d.DataValue, err = time.ParseDuration(stringValue)
		if err != nil {
			return err
		}
	case TypeTimestamp:
		d.DataValue, err = time.Parse(brokerTimeFormat, stringValue)
		if err != nil {
			return err
		}
	case TypeAddress:
		d.DataValue = net.ParseIP(stringValue)
		if d.DataValue == nil {
			return fmt.Errorf("JSON string encoded Address (%s) failed to parse", stringValue)
		}
	case TypeSubnet:
		_, d.DataValue, err = net.ParseCIDR(stringValue)
		if err != nil {
			return err
		}
	case TypePort:
		service, err := ParseService(stringValue)
		if err != nil {
			return err
		}
		d.DataValue = service
	case TypeVector:
		va, ok := v.([]interface{})
		if !ok {
			return fmt.Errorf("expected Vector type to be serialized as JSON array but got type %T value %v", v, v)
		}

		datas := make([]Data, len(va))
		for i, intf := range va {
			m, ok := intf.(map[string]interface{})
			if !ok {
				return fmt.Errorf("expected Vector type elements to be serialized as JSON objects but got type %T value %v",
					intf, intf)
			}
			err = datas[i].decode(&m)
			if err != nil {
				return fmt.Errorf("error decoding Vector element %d: %w", i, err)
			}
		}
		d.DataValue = datas
	case TypeSet:
		sa, ok := v.([]interface{})
		if !ok {
			return fmt.Errorf("expected Set type to be serialized as JSON array but got type %T value %v", v, v)
		}

		datas := make(map[Data]struct{}, len(sa))
		for i, intf := range sa {
			m, ok := intf.(map[string]interface{})
			if !ok {
				return fmt.Errorf("expected Set type elements to be serialized as JSON objects but got type %T value %v",
					intf, intf)
			}

			var dElem Data
			err = dElem.decode(&m)
			if err != nil {
				return fmt.Errorf("error decoding Set element %d: %w", i, err)
			}

			if _, ok := datas[dElem]; ok {
				return fmt.Errorf("duplicate Set element %d: %#v", i, dElem)
			}

			datas[dElem] = struct{}{}
		}
		d.DataValue = datas
	case TypeTable:
		ta, ok := v.([]interface{})
		if !ok {
			return fmt.Errorf("expected Table type to be serialized as JSON array but got type %T value %v", v, v)
		}

		datas := make(map[Data]Data, len(ta))
		for i, intf := range ta {
			m, ok := intf.(map[string]interface{})
			if !ok {
				return fmt.Errorf("expected Table type elements to be serialized as JSON objects but got type %T value %v",
					intf, intf)
			}

			if len(m) != 2 {
				return fmt.Errorf("expected Table type elements to have two properties but we have %d", len(m))
			}

			mk, ok := m["key"]
			if !ok {
				return fmt.Errorf("expected Table type elements to have a property named \"key\"")
			}

			mkm, ok := mk.(map[string]interface{})
			if !ok {
				return fmt.Errorf("expected Table type elements keys be serialized as JSON objects but got type %T value %v",
					mkm, mkm)
			}

			var dKey Data
			err = dKey.decode(&mkm)
			if err != nil {
				return fmt.Errorf("error decoding Table key element %d: %w", i, err)
			}

			mv, ok := m["value"]
			if !ok {
				return fmt.Errorf("expected Table type elements to have a property named \"value\"")
			}

			mvm, ok := mv.(map[string]interface{})
			if !ok {
				return fmt.Errorf("expected Table type elements values be serialized as JSON objects but got type %T value %v",
					mkm, mkm)
			}

			var dValue Data
			err = dValue.decode(&mvm)
			if err != nil {
				return fmt.Errorf("error decoding Table value element %d: %w", i, err)
			}

			if _, ok := datas[dKey]; ok {
				return fmt.Errorf("duplicate Table key %d: %#v", i, dKey)
			}

			datas[dKey] = dValue
		}
		d.DataValue = datas
	}

	return nil

}

// UnmarshalJSON implemnts the Unmarshaller interface for Data. It calls json.Unmarshal to produce a map[string]interface{}
// which is then passed to Data.decode() which does the heavy lifting.
func (d *Data) UnmarshalJSON(b []byte) error {
	var rawData map[string]interface{}

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()

	if err := dec.Decode(&rawData); err != nil {
		return err
	}

	return d.decode(&rawData)
}

// formatTimespan implements the string encoding of the zeek timespan type in the format specific to the broker WS API.
func formatTimespan(duration time.Duration) string {
	trim := func(s string) string {
		return strings.TrimRight(strings.TrimRight(s, "0"), ".")
	}

	ret := float64(duration) // Nanoseconds
	if duration.Abs() < time.Millisecond {
		return trim(fmt.Sprintf("%.4f", ret)) + "ns"
	}
	if duration.Abs() < time.Second {
		return trim(fmt.Sprintf("%.4f", ret/float64(time.Millisecond))) + "ms"
	}
	if duration.Abs() < time.Minute {
		return trim(fmt.Sprintf("%.4f", ret/float64(time.Second))) + "s"
	}
	if duration.Abs() < time.Hour {
		return trim(fmt.Sprintf("%.4f", ret/float64(time.Minute))) + "min"
	}
	if duration.Abs() < nanosecondsIn24Hours {
		return trim(fmt.Sprintf("%.4f", ret/float64(time.Hour))) + "h"
	}
	return trim(fmt.Sprintf("%.4f", ret/nanosecondsIn24Hours)) + "d"
}

// MarshalJSON implements the Marshaller interface for Data, taking care specific cases where json.Marshal doesn't
// produce output compliant to the zeek broker websocket JSON encoding (e.g., timestamps, ports, etc).
func (d *Data) MarshalJSON() ([]byte, error) {
	switch d.DataType {
	case TypeTimestamp:
		ts, ok := d.DataValue.(time.Time)
		if !ok {
			return nil, fmt.Errorf("expected a time.Time as the DataValue but got a %T %v", ts, ts)
		}
		return json.Marshal(map[string]interface{}{
			"@data-type": d.DataType,
			"data":       ts.Format(brokerTimeFormat),
		})
	case TypeTimespan:
		dur, ok := d.DataValue.(time.Duration)
		if !ok {
			return nil, fmt.Errorf("expected a time.Duration as the DataValue but got a %T %v", dur, dur)
		}
		return json.Marshal(map[string]interface{}{
			"@data-type": d.DataType,
			"data":       formatTimespan(dur),
		})
	case TypePort:
		serv, ok := d.DataValue.(Service)
		if !ok {
			return nil, fmt.Errorf("expected a encoding.Service as the DataValue but got a %T %v", serv, serv)
		}
		return json.Marshal(map[string]interface{}{
			"@data-type": d.DataType,
			"data":       fmt.Sprintf("%d/%s", serv.Port, serv.Protocol.String()),
		})
	default:
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)

		if err := enc.Encode(map[string]interface{}{
			"@data-type": d.DataType,
			"data":       d.DataValue,
		}); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}
}

// String implements the Stringer interface for encoding.Data and produces a compact string representation.
func (d *Data) String() string {
	return fmt.Sprintf("\"%v\": %s", d.DataValue, d.DataType.String())
}
