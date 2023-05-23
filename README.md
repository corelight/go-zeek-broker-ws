# Zeek Broker websocket interface library for Golang

This library implements the [Zeek Broker websocket interface](https://docs.zeek.org/projects/broker/en/master/web-socket.html).

Helper functions are provided with reasonable mappings between Zeek and Go types.

## Usage

The library contains two packages:

### `encoding`
`encoding` models the recursive JSON-based data structure for representing Zeek types (the `Data` struct), as well as 
the Broker encoding of Zeek events and error messages (the `DataMessage` struct). These structures can (and in some cases must)
be used directly (see the subscription example in the next section), but some helper functions make building messages
more convenient. To construct a Zeek `string` manually:

```go
zeekString := encoding.Data{
    DataType:  encoding.TypeString,
    DataValue: "foo",
}
```

...or via the helper:
```go
zeekString := encoding.String("foo")
```

The `vector` helper in particular is a variadic function that accepts `encoding.Data`:
```go
zeekVector := encoding.Vector(encoding.Count(1), encoding.Count(2), encoding.Count(3))
```

Finally, events can be created directly:
```go
zeekEvent := encoding.NewEvent("some_event_name", zeekVector, zeekString)
```

### `client`
`client` provides the websocket glue to speak to the broker WS API, wrapping `github.com/gorilla/websocket`:

To publish an event:
```go
broker, err := client.NewClient(...)

err := broker.PublishEvent("/the/topic", zeekEvent)
```

Topic subscriptions are passed as a slice of strings to `client.Newclient()`. The `ReadEvent()` method of the client 
returns a single event from Broker (on any of the subscribed topics), or an error that could occur in the library itself
ir errors received from Broker):

```go
broker, err := client.NewClient(..., []string{"/the/topic"})

topic, zeekEvent, err := broker.ReadEvent()
```

If the broker connection is closed gracefully, the `client.IsNormalWebsocketClose()` function can be used to check
the returned error.

Client code must access event argument values via type assertions:
```go
if len(evt.Arguments) < someConstantGreaterOrEqualToOne {
	// handle too few arguments case
}

if evt.Arguments[0].DataType != encoding.TypeString {
	// handle unexpected data type case
}

stringArgument0, ok := evt.Arguments[0].DataValue.(string)
if !ok {
	// handle type assertion error - this would indicate a bug (we trust+verify).
}

// now we use stringArgument0
```

More advanced handling of the websocket connection (e.g., setting timeouts, handling re-connection, etc.) is best implemented
as a wrapper of `client.Client`, or a new/replacement implementation that uses the `encoding` package (contributions/PRs are welcome!).

Similarly, asynchronous handling and dispatching of events received via subscriptions would best implemented as
a `Client.ReadEvent()` wrapper.

## Broker TLS details

Broker network connections (both native, and the websocket interface) enable TLS by default with an odd configuration
that disables host verification and selects a set of cipher that allow encryption without certificates. To use this mode
requires passing `weirdtls.BrokerDefaultTLSDialer` as the dailer function argument to `encoding.NewClient`. Note that
this pulls in OpenSSL as a dependency.

Alternatively the standard library `crypto/tls` implementation can be used if both sides (the client and zeek/broker)
is configured to use TLS with certificates. This library provides a convenient helper function 
(`securetls.MakeSecureDialer()`) that returns a dialer function given PEM files for the CA and client certificate/key. 
See [this btest case](tests/btests/receive_event_certs.test) for an example of this configuration.

Finally, TLS can be turned off for broker connections using `redef Broker::disable_ssl = T;`. 
See [this btest case](tests/btests/receive_event_nossl.test) for an example where `encoding.NewClient` is called
with arguments for insecure operation.

## Ping/pong example

Running a zeek-side broker script:
```console
$ cd example/
$ zeek listen.zeek
```

Now run the example:
```console
$ cd example/
$ go build && ./example 
2023/05/05 12:56:55 connected to remote endpoint with UUID=c9b0bfd6-3b8d-5de2-a51c-9af7b81aaad3 version=2.5.0-dev
2023/05/05 12:56:55 > topic=/topic/test | event ping("my-message": string, "1": count)
2023/05/05 12:56:55 < topic=/topic/test | event pong("my-message": string, "2": count)
2023/05/05 12:56:56 > topic=/topic/test | event ping("my-message": string, "2": count)
2023/05/05 12:56:56 < topic=/topic/test | event pong("my-message": string, "3": count)
```

...meanwhile the zeek script:
```console
peer added, [id=e537f8b4-de32-52ea-9587-4e6e15bdfe20, network=[address=127.0.0.1, bound_port=50690/tcp]]
receiver got ping: my-message, 1
receiver got ping: my-message, 2
```

## Dev workflow

Most helpers and basic data structures in the `encoding` package have unit tests (run `go test ./...`).

End-to-end tests are implemented as two [btest](https://github.com/zeek/btest) cases (run `cd tests/; btest`).
