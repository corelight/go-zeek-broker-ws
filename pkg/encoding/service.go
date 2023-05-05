// Copyright (c) 2023, Corelight, Inc. All rights reserved.

package encoding

import (
	"fmt"
	"strconv"
	"strings"
)

//go:generate go-enum --marshal

// Protocol represents an L4 protocol
/*
ENUM(
TCP = "tcp"
UDP = "udp"
ICMP = "icmp"
Unknown = "?"
)
*/
type Protocol string

type Service struct {
	Port     uint16
	Protocol Protocol
}

const numServiceParts = 2

func ParseService(stringValue string) (Service, error) {
	parts := strings.Split(stringValue, "/")
	if len(parts) != numServiceParts {
		return Service{}, fmt.Errorf("Port value of %s has too many (%d) parts", stringValue, len(parts))
	}
	p, err := strconv.ParseUint(parts[0], 10, 16)
	if err != nil {
		return Service{}, err
	}
	l4, err := ParseProtocol(parts[1])
	if err != nil {
		return Service{}, err
	}

	return Service{
		Port:     uint16(p),
		Protocol: l4,
	}, nil
}

func (s *Service) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%d/%s\"", s.Port, s.Protocol.String())), nil
}
