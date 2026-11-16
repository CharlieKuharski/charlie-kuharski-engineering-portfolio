package nats

import "context"

// EventHandler is the only domain type needed by this adapter.
type EventHandler func(ctx context.Context, payload []byte) error
