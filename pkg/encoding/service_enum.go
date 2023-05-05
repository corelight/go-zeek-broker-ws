// Copyright (c) 2023, Corelight, Inc. All rights reserved.

package encoding

import (
	"errors"
	"fmt"
)

const (
	// ProtocolTCP is a Protocol of type TCP.
	ProtocolTCP Protocol = "tcp"
	// ProtocolUDP is a Protocol of type UDP.
	ProtocolUDP Protocol = "udp"
	// ProtocolICMP is a Protocol of type ICMP.
	ProtocolICMP Protocol = "icmp"
	// ProtocolUnknown is a Protocol of type Unknown.
	ProtocolUnknown Protocol = "?"
)

var ErrInvalidProtocol = errors.New("not a valid Protocol")

// String implements the Stringer interface.
func (x Protocol) String() string {
	return string(x)
}

// String implements the Stringer interface.
func (x Protocol) IsValid() bool {
	_, err := ParseProtocol(string(x))
	return err == nil
}

var _ProtocolValue = map[string]Protocol{
	"tcp":  ProtocolTCP,
	"udp":  ProtocolUDP,
	"icmp": ProtocolICMP,
	"?":    ProtocolUnknown,
}

// ParseProtocol attempts to convert a string to a Protocol.
func ParseProtocol(name string) (Protocol, error) {
	if x, ok := _ProtocolValue[name]; ok {
		return x, nil
	}
	return Protocol(""), fmt.Errorf("%s is %w", name, ErrInvalidProtocol)
}

// MarshalText implements the text marshaller method.
func (x Protocol) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Protocol) UnmarshalText(text []byte) error {
	tmp, err := ParseProtocol(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
