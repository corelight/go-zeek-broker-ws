// Copyright (c) 2023, Corelight, Inc. All rights reserved.

package encoding

import (
	"math"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestData_Bool(t *testing.T) {
	wantData := Data{
		DataType:  "boolean",
		DataValue: true,
	}

	gotData := Boolean(true)

	if !reflect.DeepEqual(wantData, gotData) {
		t.Errorf("output value incorrect, wanted: \n\t%#v\ngot: \n\t%#v", wantData, gotData)
	}
}

func TestData_Count(t *testing.T) {
	wantData := Data{
		DataType:  "count",
		DataValue: uint64(0),
	}

	gotData := Count(0)

	if !reflect.DeepEqual(wantData, gotData) {
		t.Errorf("output value incorrect, wanted: \n\t%#v\ngot: \n\t%#v", wantData, gotData)
	}

	wantData = Data{
		DataType:  "count",
		DataValue: uint64(math.MaxUint64),
	}

	gotData = Count(math.MaxUint64)

	if !reflect.DeepEqual(wantData, gotData) {
		t.Errorf("output value incorrect, wanted: \n\t%#v\n\t%#v", wantData, gotData)
	}
}

func TestData_Integer(t *testing.T) {
	wantData := Data{
		DataType:  "integer",
		DataValue: int64(math.MaxInt64),
	}

	gotData := Integer(math.MaxInt64)

	if !reflect.DeepEqual(wantData, gotData) {
		t.Errorf("output value incorrect, wanted: \n\t%#v\n\t%#v", wantData, gotData)
	}

	wantData = Data{
		DataType:  "integer",
		DataValue: int64(math.MinInt64),
	}

	gotData = Integer(math.MinInt64)

	if !reflect.DeepEqual(wantData, gotData) {
		t.Errorf("output value incorrect, wanted: \n\t%#v\n\t%#v", wantData, gotData)
	}
}

func TestData_Real(t *testing.T) {
	wantData := Data{
		DataType:  "real",
		DataValue: math.Pi,
	}

	gotData := Real(math.Pi)

	if !reflect.DeepEqual(wantData, gotData) {
		t.Errorf("output value incorrect, wanted: \n\t%#v\ngot: \n\t%#v", wantData, gotData)
	}
}

func TestData_Timespan(t *testing.T) {
	wantData := Data{
		DataType:  "timespan",
		DataValue: time.Second,
	}

	gotData := Timespan(time.Second)

	if !reflect.DeepEqual(wantData, gotData) {
		t.Errorf("output value incorrect, wanted: \n\t%#v\ngot: \n\t%#v", wantData, gotData)
	}
}

func TestData_Timestamp(t *testing.T) {
	wantData := Data{
		DataType:  "timestamp",
		DataValue: time.Date(2006, time.January, 2, 15, 4, 5, 999000000, time.UTC),
	}

	timeString := "2006-01-02T15:04:05.999"
	ts, err := time.Parse("2006-01-02T15:04:05.999", timeString)
	if err != nil {
		t.Fatal(err)
	}

	gotData := Timestamp(ts)

	if !reflect.DeepEqual(wantData, gotData) {
		t.Errorf("output value incorrect, wanted: \n\t%#v\ngot: \n\t%#v", wantData, gotData)
	}
}

func TestData_String(t *testing.T) {
	wantData := Data{
		DataType:  "string",
		DataValue: "foo",
	}

	gotData := String("foo")

	if !reflect.DeepEqual(wantData, gotData) {
		t.Errorf("output value incorrect, wanted: \n\t%#v\ngot: \n\t%#v", wantData, gotData)
	}
}

func TestData_EnumValue(t *testing.T) {
	wantData := Data{
		DataType:  "enum-value",
		DataValue: "foo",
	}

	gotData := EnumValue("foo")

	if !reflect.DeepEqual(wantData, gotData) {
		t.Errorf("output value incorrect, wanted: \n\t%#v\ngot: \n\t%#v", wantData, gotData)
	}
}

func TestData_Subnet(t *testing.T) {
	_, network, err := net.ParseCIDR("1.2.3.0/24")
	if err != nil {
		t.Fatal(err)
	}

	wantData := Data{
		DataType:  "subnet",
		DataValue: "1.2.3.0/24",
	}

	gotData := Subnet(*network)

	if !reflect.DeepEqual(wantData, gotData) {
		t.Errorf("output value incorrect, wanted: \n\t%#v\ngot: \n\t%#v", wantData, gotData)
	}
}

func TestData_Addr(t *testing.T) {
	addr := net.ParseIP("1.2.3.4")

	wantData := Data{
		DataType:  "address",
		DataValue: "1.2.3.4",
	}

	gotData := Address(addr)

	if !reflect.DeepEqual(wantData, gotData) {
		t.Errorf("output value incorrect, wanted: \n\t%#v\ngot: \n\t%#v", wantData, gotData)
	}
}

func TestData_Port(t *testing.T) {
	service := Service{
		Port:     443,
		Protocol: ProtocolTCP,
	}

	wantData := Data{
		DataType:  "port",
		DataValue: service,
	}

	gotData := Port(service)

	if !reflect.DeepEqual(wantData, gotData) {
		t.Errorf("output value incorrect, wanted: \n\t%#v\ngot: \n\t%#v", wantData, gotData)
	}
}

func TestData_Set(t *testing.T) {
	// Note: two attempts because maps are unordered, this seems
	//       silly, but less silly than sorting the output.
	var wantData1 = Data{
		DataType: "set",
		DataValue: []Data{
			{DataType: "string", DataValue: "bar"},
			{DataType: "string", DataValue: "foo"},
		},
	}

	var wantData2 = Data{
		DataType: "set",
		DataValue: []Data{
			{DataType: "string", DataValue: "foo"},
			{DataType: "string", DataValue: "bar"},
		},
	}

	gotData := Set(map[Data]struct{}{
		String("foo"): {},
		String("bar"): {},
	})

	if !reflect.DeepEqual(wantData1, gotData) {
		if !reflect.DeepEqual(wantData2, gotData) {
			t.Errorf("output value incorrect, wanted either: \n\t%#v\nor \n\t%#v\ngot: \n\t%#v",
				wantData1, wantData2, gotData)
		}
	}
}

func TestData_Table(t *testing.T) {
	// Note: two attempts because maps are unordered, this seems
	//       silly, but less silly than sorting the output.
	var wantData1 = Data{
		DataType: "table",
		DataValue: []map[string]Data{
			{
				"key":   String("foo"),
				"value": String("bar"),
			},
			{
				"key":   String("bar"),
				"value": String("foo"),
			},
		},
	}

	var wantData2 = Data{
		DataType: "table",
		DataValue: []map[string]Data{
			{
				"key":   String("bar"),
				"value": String("foo"),
			},
			{
				"key":   String("foo"),
				"value": String("bar"),
			},
		},
	}

	gotData := Table(map[Data]Data{
		String("foo"): String("bar"),
		String("bar"): String("foo"),
	})

	if !reflect.DeepEqual(wantData1, gotData) {
		if !reflect.DeepEqual(wantData2, gotData) {
			t.Errorf("output value incorrect, wanted either: \n\t%#v\nor \n\t%#v\ngot: \n\t%#v",
				wantData1, wantData2, gotData)
		}
	}
}

func TestData_Vector(t *testing.T) {
	var wantData = Data{
		DataType: "vector",
		DataValue: []Data{
			{DataType: "string", DataValue: "foo"},
			{DataType: "string", DataValue: "bar"},
		},
	}

	gotData := Vector(
		String("foo"),
		String("bar"),
	)

	if !reflect.DeepEqual(wantData, gotData) {
		t.Errorf("output value incorrect, wanted: \n\t%#v\ngot: \n\t%#v", wantData, gotData)
	}
}

func TestData_Event(t *testing.T) {
	var wantData = Vector(
		Count(1),
		Count(1),
		Vector(
			String("test_event"),
			Vector(
				String("foo"),
				String("bar"),
			),
		),
	)

	var wantEvent = DataMessage{
		ConstType: "data-message",
		Topic:     "test_topic",
		Data:      &wantData,
	}

	gotEvent := NewEvent(
		"test_event",
		String("foo"),
		String("bar"),
	).Encode("test_topic")

	if wantEvent.Topic != gotEvent.Topic || wantEvent.ConstType != gotEvent.ConstType {
		t.Fatal("unit test is b0rk")
	}

	if !reflect.DeepEqual(wantEvent.Data, gotEvent.Data) {
		t.Errorf("output value incorrect, wanted: \n\t%#v\ngot: \n\t%#v", wantEvent.Data, gotEvent.Data)
	}
}

func TestData_None(t *testing.T) {
	wantData := Data{
		DataType:  "none",
		DataValue: nil,
	}

	gotData := None()

	if !reflect.DeepEqual(wantData, gotData) {
		t.Errorf("output value incorrect, wanted: \n\t%#v\ngot: \n\t%#v", wantData, gotData)
	}
}
