package types

import "time"

type Message struct {
	// The origin URL of the websocket message.
	Origin string

	// The raw content of the websocket message.
	WebsocketMessage []byte

	Time time.Time
}
