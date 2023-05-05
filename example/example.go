// Copyright (c) 2023, Corelight, Inc. All rights reserved.

package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/corelight/go-zeek-broker-ws/pkg/client"
	"github.com/corelight/go-zeek-broker-ws/pkg/encoding"
	"github.com/gorilla/websocket"
)

func main() {
	topic := "/topic/test"
	event := "ping"

	ctx, cancel := context.WithCancel(context.Background())

	broker, err := client.NewClient(ctx, "localhost:9997", false, []string{topic})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Printf("shutting down websocket")
		_ = broker.Close()
	}()

	uuid, version := broker.RemoteEndpointInfo()
	log.Printf("connected to remote endpoint with UUID=%s version=%s\n", uuid, version)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
		log.Printf("received signal, shutting down websocket")
		_ = broker.Close()
	}()

	i := uint64(1)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			message := encoding.String("my-message")
			count := encoding.Count(i)
			evt := encoding.NewEvent(event, message, count)

			log.Printf("> topic=%s | %s\n", topic, evt)

			err = broker.PublishEvent(topic, evt)
			if err != nil {
				if errors.Is(err, websocket.ErrCloseSent) {
					log.Printf("broker connection closed while publishing event")
					return
				}
				log.Fatal(err)
			}

			rcvTopic, rcvEvt, err := broker.ReadEvent()
			if err != nil {
				if client.IsNormalWebsocketClose(err) {
					log.Printf("broker connection closed while reading event")
					return
				}
				log.Fatal(err)
			}

			log.Printf("< topic=%s | %s\n", rcvTopic, rcvEvt)

			i++

			time.Sleep(time.Second)
		}
	}
}
