// Copyright (c) 2023, Corelight, Inc. All rights reserved.

package encoding

import (
	"net"
	"reflect"
	"testing"
	"time"
)

//nolint:gocognit // non-complex repetition of subtests
func TestData_UnmarshalJSON(t *testing.T) {
	var dur time.Duration
	ipv6 := net.ParseIP("2001:db8::")
	if ipv6 == nil {
		t.Fatal("test data is invalid")
	}
	ipv4 := net.ParseIP("196.25.1.1")
	if ipv4 == nil {
		t.Fatal("test data is invalid")
	}
	_, ipv6Subnet, err := net.ParseCIDR("2001:db8::/127")
	if err != nil {
		t.Fatal(err)
	}
	_, ipv4Subnet, err := net.ParseCIDR("196.25.0.0/16")
	if err != nil {
		t.Fatal(err)
	}
	timeString := "2006-01-02T15:04:05.999"
	ts, err := time.Parse("2006-01-02T15:04:05.999", timeString)
	if err != nil {
		t.Fatal(err)
	}
	serv := Service{
		Port:     25,
		Protocol: "tcp",
	}
	vec := []Data{
		{
			DataType:  "count",
			DataValue: uint64(42),
		},
		{
			DataType:  "integer",
			DataValue: int64(23),
		},
	}
	set := map[Data]struct{}{
		{
			DataType:  "string",
			DataValue: "foo",
		}: {},
		{
			DataType:  "string",
			DataValue: "bar",
		}: {},
	}
	table := map[Data]Data{
		{
			DataType:  "string",
			DataValue: "first-name",
		}: {
			DataType:  "string",
			DataValue: "John",
		},
		{
			DataType:  "string",
			DataValue: "last-name",
		}: {
			DataType:  "string",
			DataValue: "Doe",
		},
	}

	tests := []struct {
		name     string
		arg      []byte
		want     Data
		wantType reflect.Kind
		wantErr  bool
	}{
		{name: "count valid", want: Data{DataType: TypeCount, DataValue: uint64(123)},
			wantType: reflect.Uint64, wantErr: false, arg: []byte(`
		{
		  "@data-type": "count",
		  "data": 123
		}
		`)},
		{name: "count invalid value sign", want: Data{DataType: TypeCount, DataValue: uint64(123)},
			wantType: reflect.Uint64, wantErr: true, arg: []byte(`
		{
		  "@data-type": "count",
		  "data": -123
		}
		`)},
		{name: "count invalid value type", want: Data{DataType: TypeCount, DataValue: uint64(123)},
			wantType: reflect.Uint64, wantErr: true, arg: []byte(`
		{
		  "@data-type": "count",
		  "data": "123""
		}
		`)},
		{name: "integer valid", want: Data{DataType: TypeInteger, DataValue: int64(-7)},
			wantType: reflect.Int64, wantErr: false, arg: []byte(`
		{
		  "@data-type": "integer",
		  "data": -7
		}
		`)},
		{name: "integer invalid range 1", want: Data{DataType: TypeInteger, DataValue: int64(-7)},
			wantType: reflect.Int64, wantErr: true, arg: []byte(`
		{
		  "@data-type": "integer",
		  "data": 18446744073709551615
		}
		`)},
		{name: "integer invalid range 2", want: Data{DataType: TypeInteger, DataValue: int64(-7)},
			wantType: reflect.Int64, wantErr: true, arg: []byte(`
		{
		  "@data-type": "integer",
		  "data": -18446744073709551615
		}
		`)},
		{name: "timespan valid 1ns", want: Data{DataType: TypeTimespan, DataValue: time.Nanosecond},
			wantType: reflect.TypeOf(dur).Kind(), wantErr: false, arg: []byte(`
		{
		  "@data-type": "timespan",
		  "data": "1ns"
		}
		`)},
		{name: "timespan valid 1ms", want: Data{DataType: TypeTimespan, DataValue: time.Millisecond},
			wantType: reflect.TypeOf(dur).Kind(), wantErr: false, arg: []byte(`
		{
		  "@data-type": "timespan",
		  "data": "1ms"
		}
		`)},
		{name: "timespan valid 1s", want: Data{DataType: TypeTimespan, DataValue: time.Second},
			wantType: reflect.TypeOf(dur).Kind(), wantErr: false, arg: []byte(`
		{
		  "@data-type": "timespan",
		  "data": "1s"
		}
		`)},
		{name: "timespan valid 1min", want: Data{DataType: TypeTimespan, DataValue: time.Minute},
			wantType: reflect.TypeOf(dur).Kind(), wantErr: false, arg: []byte(`
		{
		  "@data-type": "timespan",
		  "data": "1min"
		}
		`)},
		{name: "timespan valid 1h", want: Data{DataType: TypeTimespan, DataValue: time.Hour},
			wantType: reflect.TypeOf(dur).Kind(), wantErr: false, arg: []byte(`
		{
		  "@data-type": "timespan",
		  "data": "1h"
		}
		`)},
		{name: "timespan valid 1d", want: Data{DataType: TypeTimespan, DataValue: time.Hour * 24},
			wantType: reflect.TypeOf(dur).Kind(), wantErr: false, arg: []byte(`
		{
		  "@data-type": "timespan",
		  "data": "1d"
		}
		`)},
		{name: "timespan valid 1.234567d", want: Data{DataType: TypeTimespan, DataValue: time.Nanosecond * 106666588800000},
			wantType: reflect.TypeOf(dur).Kind(), wantErr: false, arg: []byte(`
		{
		  "@data-type": "timespan",
		  "data": "1.234567d"
		}
		`)},
		{name: "timestamp valid", want: Data{DataType: TypeTimestamp, DataValue: ts},
			wantType: reflect.TypeOf(ts).Kind(), wantErr: false, arg: []byte(`
		{
		  "@data-type": "timestamp",
		  "data": "2006-01-02T15:04:05.999"
		}
		`)},
		{name: "address IPv6 valid", want: Data{DataType: TypeAddress, DataValue: ipv6},
			wantType: reflect.TypeOf(ipv6).Kind(), wantErr: false, arg: []byte(`
		{
		  "@data-type": "address",
		  "data": "2001:db8::"
		}
		`)},
		{name: "address IPv4 valid", want: Data{DataType: TypeAddress, DataValue: ipv4},
			wantType: reflect.TypeOf(ipv4).Kind(), wantErr: false, arg: []byte(`
		{
		  "@data-type": "address",
		  "data": "196.25.1.1"
		}
		`)},
		{name: "subnet IPv6 valid", want: Data{DataType: TypeSubnet, DataValue: ipv6Subnet},
			wantType: reflect.TypeOf(ipv6Subnet).Kind(), wantErr: false, arg: []byte(`
		{
		  "@data-type": "subnet",
		  "data": "2001:db8::/127"
		}
		`)},
		{name: "subnet IPv4 valid", want: Data{DataType: TypeSubnet, DataValue: ipv4Subnet},
			wantType: reflect.TypeOf(ipv4Subnet).Kind(), wantErr: false, arg: []byte(`
		{
		  "@data-type": "subnet",
		  "data": "196.25.0.0/16"
		}
		`)},
		{name: "service valid", want: Data{DataType: TypePort, DataValue: serv},
			wantType: reflect.TypeOf(serv).Kind(), wantErr: false, arg: []byte(`
		{
		  "@data-type": "port",
		  "data": "25/tcp"
		}
		`)},
		{name: "service invalid everything", want: Data{DataType: TypePort, DataValue: serv},
			wantType: reflect.TypeOf(serv).Kind(), wantErr: true, arg: []byte(`
		{
		  "@data-type": "port",
		  "data": "sdfsdf"
		}
		`)},
		{name: "service invalid port", want: Data{DataType: TypePort, DataValue: serv},
			wantType: reflect.TypeOf(serv).Kind(), wantErr: true, arg: []byte(`
		{
		  "@data-type": "port",
		  "data": "sdfsdf/tcp"
		}
		`)},
		{name: "service invalid protocol", want: Data{DataType: TypePort, DataValue: serv},
			wantType: reflect.TypeOf(serv).Kind(), wantErr: true, arg: []byte(`
		{
		  "@data-type": "port",
		  "data": "25/sctp"
		}
		`)},
		{name: "vector valid", want: Data{DataType: TypeVector, DataValue: vec},
			wantType: reflect.TypeOf(vec).Kind(), wantErr: false, arg: []byte(`
		{
		  "@data-type": "vector",
		  "data": [
			  {
				"@data-type": "count",
				"data": 42
			  },
			  {
				"@data-type": "integer",
				"data": 23
			  }
			]
		}
		`)},
		{name: "set valid", want: Data{DataType: TypeSet, DataValue: set},
			wantType: reflect.TypeOf(set).Kind(), wantErr: false, arg: []byte(`
		{
			"@data-type": "set",
			"data": [
			  {
				"@data-type": "string",
				"data": "foo"
			  },
			  {
				"@data-type": "string",
				"data": "bar"
			  }
			]
		}
		`)},
		{name: "set invalid", want: Data{DataType: TypeSet, DataValue: set},
			wantType: reflect.TypeOf(set).Kind(), wantErr: true, arg: []byte(`
		{
			"@data-type": "set",
			"data": [
			  {
				"@data-type": "string",
				"data": "bar"
			  },
			  {
				"@data-type": "string",
				"data": "bar"
			  }
			]
		}
		`)},
		{name: "table valid", want: Data{DataType: TypeTable, DataValue: table},
			wantType: reflect.TypeOf(table).Kind(), wantErr: false, arg: []byte(`
		{
		  "@data-type": "table",
		  "data": [
			{
			  "key": {
				"@data-type": "string",
				"data": "first-name"
			  },
			  "value": {
				"@data-type": "string",
				"data": "John"
			  }
			},
			{
			  "key": {
				"@data-type": "string",
				"data": "last-name"
			  },
			  "value": {
				"@data-type": "string",
				"data": "Doe"
			  }
			}
		  ]
		}
		`)},
		{name: "table invalid", want: Data{DataType: TypeTable, DataValue: table},
			wantType: reflect.TypeOf(table).Kind(), wantErr: true, arg: []byte(`
		{
		  "@data-type": "table",
		  "data": [
			{
			  "key": {
				"@data-type": "string",
				"data": "first-name"
			  },
			  "value": {
				"@data-type": "string",
				"data": "John"
			  }
			},
			{
			  "key": {
				"@data-type": "string",
				"data": "first-name"
			  },
			  "value": {
				"@data-type": "string",
				"data": "Doe"
			  }
			}
		  ]
		}
		`)},
		{name: "none valid", want: Data{DataType: TypeNone, DataValue: nil},
			wantType: reflect.Invalid, wantErr: false, arg: []byte(`
		{
		  "@data-type": "none",
		  "data": {}
		}`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Data
			err := d.UnmarshalJSON(tt.arg)
			if !tt.wantErr && (err != nil) {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr && (err == nil) {
				if d.DataValue == nil && d.DataType == TypeNone {
					return
				}
				kind := reflect.TypeOf(d.DataValue).Kind()
				if kind != tt.wantType {
					t.Errorf("output type incorrect, want: %#v, got: %#v", tt.wantType, kind)
				}
				if !reflect.DeepEqual(tt.want, d) {
					t.Errorf("output value incorrect, want: %#v, got: %#v", tt.want, d)
				}
			}
		})
	}
}

func Test_formatTimespan(t *testing.T) {
	type args struct {
		duration time.Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "nanoseconds", args: args{duration: time.Nanosecond * 10}, want: "10ns"},
		{name: "milliseconds", args: args{duration: time.Nanosecond * 1500000}, want: "1.5ms"},
		{name: "seconds", args: args{duration: time.Nanosecond * 1500000000}, want: "1.5s"},
		{name: "minutes", args: args{duration: time.Nanosecond * 90000000000}, want: "1.5min"},
		{name: "hours", args: args{duration: time.Nanosecond * 5400000000000}, want: "1.5h"},
		{name: "days", args: args{duration: time.Nanosecond * 129600000000000}, want: "1.5d"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatTimespan(tt.args.duration); got != tt.want {
				t.Errorf("formatTimespan() = %v, want %v", got, tt.want)
			}
		})
	}
}
