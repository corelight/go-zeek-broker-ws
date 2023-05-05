// Copyright (c) 2023, Corelight, Inc. All rights reserved.

package encoding

import (
	"errors"
	"fmt"
)

const (
	// TypeBoolean is a Type of type Boolean.
	// Native JSON boolean (maps to bool)
	TypeBoolean Type = "boolean"
	// TypeCount is a Type of type Count.
	// 64 bit unsigned integer (maps to uint64)
	TypeCount Type = "count"
	// TypeInteger is a Type of type Integer.
	// Native JSON (signed) integer (maps to int64)
	TypeInteger Type = "integer"
	// TypeReal is a Type of type Real.
	// Native JSON "number" type (maps to float64)
	TypeReal Type = "real"
	// TypeTimespan is a Type of type Timespan.
	// String-encoded time span (maps to time.Duration)
	TypeTimespan Type = "timespan"
	// TypeTimestamp is a Type of type Timestamp.
	// ISO 8601 encoded time in YYYY-MM-DDThh:mm:ss.sss format (maps to time.Time)
	TypeTimestamp Type = "timestamp"
	// TypeString is a Type of type String.
	// Native JSON string (maps to string)
	TypeString Type = "string"
	// TypeEnumValue is a Type of type EnumValue.
	// Zeek enum value mapped to native JSON string (maps to string)
	TypeEnumValue Type = "enum-value"
	// TypeAddress is a Type of type Address.
	// String-encoded IPv4/IPv6 address (maps to net.Addr)
	TypeAddress Type = "address"
	// TypeSubnet is a Type of type Subnet.
	// String-encoded IPv4/IPv6 subnet in <address>/<prefix-length> format (maps to net.IPNet)
	TypeSubnet Type = "subnet"
	// TypePort is a Type of type Port.
	// String-encoded service port in <port>/<protocol> format (maps to encoding.Service)
	TypePort Type = "port"
	// TypeVector is a Type of type Vector.
	// Sequence of encoding.Data (maps to []Data)
	TypeVector Type = "vector"
	// TypeSet is a Type of type Set.
	// Sequence of encoding.Data with distinct objects (maps to map[Data]struct{})
	TypeSet Type = "set"
	// TypeTable is a Type of type Table.
	// Map of encoding.Data keys to encoding.Data values (maps to map[Data]Data)
	TypeTable Type = "table"
	// TypeNone is a Type of type None.
	// JSON empty object, maps to nil
	TypeNone Type = "none"
)

var ErrInvalidType = errors.New("not a valid Type")

// String implements the Stringer interface.
func (x Type) String() string {
	return string(x)
}

// String implements the Stringer interface.
func (x Type) IsValid() bool {
	_, err := ParseType(string(x))
	return err == nil
}

var _TypeValue = map[string]Type{
	"boolean":    TypeBoolean,
	"count":      TypeCount,
	"integer":    TypeInteger,
	"real":       TypeReal,
	"timespan":   TypeTimespan,
	"timestamp":  TypeTimestamp,
	"string":     TypeString,
	"enum-value": TypeEnumValue,
	"address":    TypeAddress,
	"subnet":     TypeSubnet,
	"port":       TypePort,
	"vector":     TypeVector,
	"set":        TypeSet,
	"table":      TypeTable,
	"none":       TypeNone,
}

// ParseType attempts to convert a string to a Type.
func ParseType(name string) (Type, error) {
	if x, ok := _TypeValue[name]; ok {
		return x, nil
	}
	return Type(""), fmt.Errorf("%s is %w", name, ErrInvalidType)
}

// MarshalText implements the text marshaller method.
func (x Type) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Type) UnmarshalText(text []byte) error {
	tmp, err := ParseType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
