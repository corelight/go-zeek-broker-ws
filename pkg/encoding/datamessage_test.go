// Copyright (c) 2023, Corelight, Inc. All rights reserved.

package encoding

import (
	"encoding/json"
	"testing"
)

func TestDataMessage_UnmarshalJSON(t *testing.T) {
	raw := `
{
  "type": "data-message",
  "topic": "/topic/test",
  "@data-type": "vector",
  "data": [
    {
      "@data-type": "count",
      "data": 1
    },
    {
      "@data-type": "count",
      "data": 1
    },
    {
      "@data-type": "vector",
      "data": [
        {
          "@data-type": "string",
          "data": "pong"
        },
        {
          "@data-type": "vector",
          "data": [
            {
              "@data-type": "string",
              "data": "my-message"
            },
            {
              "@data-type": "count",
              "data": 2
            }
          ]
        }
      ]
    }
  ]
}
`
	var dm DataMessage
	err := json.Unmarshal([]byte(raw), &dm)
	if err != nil {
		t.Error(err)
	}

	if dm.ConstType != "data-message" {
		t.Fatalf("dm.ConstType is wrong: %+v", dm)
	}

	if dm.Topic != "/topic/test" {
		t.Fatalf("dm.Topic is wrong: %+v", dm)
	}
}
