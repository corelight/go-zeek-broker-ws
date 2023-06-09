// Copyright (C) 2023 Corelight Inc - All Rights Reserved

package client

import (
	"context"
	"net"

	"github.com/corelight/go-zeek-broker-ws/pkg/encoding"
	"github.com/gorilla/websocket"
)

type EventHandler func(topic string, event encoding.Event)

type ErrorHandler func(err error)

// AsyncSubscription runs the message handling loop given an EventHandler and optional ErrorHandler.
//
//nolint:gocognit // neccessary nesting
func AsyncSubscription(ctx context.Context, broker *Client, hm EventHandler, eh ErrorHandler) {
	if hm == nil {
		panic("Client.Handle must be passed a non-nil EventHandler")
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				topic, evt, err := broker.ReadEvent()

				if err != nil {
					e, ok := err.(*websocket.CloseError) //nolint:errorlint //oh shush
					if ok {
						// Normal EOF close
						if e.Code == websocketNormalEOFCode {
							return
						}

						// Abnormal websocket error, pass to handler then exit
						eh(err)
						return
					}

					//nolint:errorlint //oh shush
					if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
						return
					}
					eh(err)
					continue
				}

				hm(topic, evt)
			}
		}
	}()
}
