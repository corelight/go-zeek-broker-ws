// Copyright (c) 2023, Corelight, Inc. All rights reserved.

package encoding

import (
	"net"
	"time"
)

// Boolean creates an encoding.Data of boolean type given the provided bool value.
func Boolean(value bool) Data {
	return Data{
		DataType:  TypeBoolean,
		DataValue: value,
	}
}

// Count creates an encoding.Data of count type given the provided uint64 value.
func Count(value uint64) Data {
	return Data{
		DataType:  TypeCount,
		DataValue: value,
	}
}

// Integer creates an encoding.Data of integer type given the provided int64 value.
func Integer(value int64) Data {
	return Data{
		DataType:  TypeInteger,
		DataValue: value,
	}
}

// Real creates an encoding.Data of real type given the provided float64 value.
func Real(value float64) Data {
	return Data{
		DataType:  TypeReal,
		DataValue: value,
	}
}

// Timespan creates an encoding.Data of timespan type given the provided time.Duration value.
func Timespan(value time.Duration) Data {
	return Data{
		DataType:  TypeTimespan,
		DataValue: value,
	}
}

// Timespan creates an encoding.Data of timestamp type given the provided time.Time value.
func Timestamp(value time.Time) Data {
	return Data{
		DataType:  TypeTimestamp,
		DataValue: value,
	}
}

// String creates an encoding.Data of string type given the provided string value.
func String(value string) Data {
	return Data{
		DataType:  TypeString,
		DataValue: value,
	}
}

// EnumValue creates an encoding.Data of enum-value type given the provided string value.
func EnumValue(value string) Data {
	return Data{
		DataType:  TypeEnumValue,
		DataValue: value,
	}
}

// Address creates an encoding.Data of address type given the provided net.IP value.
func Address(value net.IP) Data {
	return Data{
		DataType:  TypeAddress,
		DataValue: value.String(),
	}
}

// Subnet creates an encoding.Data of subnet type given the provided net.IPNet value.
func Subnet(value net.IPNet) Data {
	return Data{
		DataType:  TypeSubnet,
		DataValue: value.String(),
	}
}

// Port creates an encoding.Data of port type given the provided encoding.Service value.
func Port(value Service) Data {
	return Data{
		DataType:  TypePort,
		DataValue: value,
	}
}

// Vector creates an encoding.Data of vector type given the provided encoding.Data values.
func Vector(elements ...Data) Data {
	return Data{
		DataType:  TypeVector,
		DataValue: elements,
	}
}

// Set creates an encoding.Data of set type given the provided map of Data to struct{} value.
func Set(value map[Data]struct{}) Data {
	valueList := make([]Data, len(value))
	i := 0
	for k := range value {
		valueList[i] = k
		i++
	}
	return Data{
		DataType:  TypeSet,
		DataValue: valueList,
	}
}

// Table creates an encoding.Data of table type given the provided map of Data to Data value.
func Table(value map[Data]Data) Data {
	kvList := make([]map[string]Data, len(value))
	i := 0
	for k, v := range value {
		kvList[i] = map[string]Data{
			"key":   k,
			"value": v,
		}
		i++
	}
	return Data{
		DataType:  TypeTable,
		DataValue: kvList,
	}
}

// None creates an encoding.Data of none type.
func None() Data {
	return Data{
		DataType:  TypeNone,
		DataValue: make(map[string]interface{}),
	}
}
