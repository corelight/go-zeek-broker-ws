// Copyright (c) 2023, Corelight, Inc. All rights reserved.

package client

import (
	"context"
	"fmt"

	"github.com/corelight/go-zeek-broker-ws/pkg/encoding"
	"github.com/gorilla/websocket"
)

// Client is a basic websocket client for publishing/subscribing to events via the Zeek broker websocket API.
type Client struct {
	conn            *websocket.Conn
	topics          []string
	ctx             context.Context
	endpointUUID    string
	endpointVersion string
}

const websocketNormalEOFCode = 1000

// IsNormalWebsocketClose returns true if err indicates a normal EOF close of the websocket.
func IsNormalWebsocketClose(err error) bool {
	e, ok := err.(*websocket.CloseError) //nolint:errorlint //oh shush
	if !ok {
		return false
	}
	// Normal EOF close
	if e.Code != websocketNormalEOFCode {
		return false
	}

	return true
}

// NewClient constructs a new websocket client to connect to the endpoint specified,
// subscribing to topics (which may be an empty list).
func NewClient(ctx context.Context, hostPort string, secure bool, topics []string) (*Client, error) {
	scheme := "ws"
	if secure {
		scheme = "wss"
	}

	url := fmt.Sprintf("%s://%s/v1/messages/json", scheme, hostPort)

	c, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	err = c.WriteJSON(topics)
	if err != nil {
		return nil, err
	}

	var ack encoding.AckMessage
	err = c.ReadJSON(&ack)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:            c,
		topics:          topics,
		ctx:             ctx,
		endpointUUID:    ack.EndpointUUID,
		endpointVersion: ack.EndpointUUID,
	}, nil
}

// ReadEvent reads a single event from broker, and returns the topic and event, or an error (including
// errors received from broker itself). The Client instance must be created with the list topic subscriptions.
func (c *Client) ReadEvent() (topic string, evt encoding.Event, retErr error) {
	var msg encoding.DataMessage
	if err := c.conn.ReadJSON(&msg); err != nil {
		return "", encoding.Event{}, err
	}

	return msg.GetEvent()
}

// PublishEvent publishes an event to the topic provided.
func (c *Client) PublishEvent(topic string, evt encoding.Event) error {
	return c.conn.WriteJSON(evt.Encode(topic))
}

// RemoteEndpointInfo returns the broker remote endpoint UUID and version received in the initial
// handshake when the websocket connection is established.
func (c *Client) RemoteEndpointInfo() (uuid string, version string) {
	return c.endpointUUID, c.endpointVersion
}

// Close closes the underlying websocket connection.
func (c *Client) Close() error {
	return c.conn.Close()
}
