package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// Bus is the extracted Core NATS adapter; durable execution is outside this sample.
type Bus struct {
	nc *nats.Conn
}

func NewBus(url string) (*Bus, error) {
	if url == "" {
		url = nats.DefaultURL
	}
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	return &Bus{nc: nc}, nil
}

func (b *Bus) Publish(ctx context.Context, topic string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}
	log.Printf("[NATS] Published to %s", topic)
	return b.nc.Publish(topic, data)
}

func (b *Bus) Subscribe(topic string, handler EventHandler) error {
	_, err := b.nc.Subscribe(topic, func(m *nats.Msg) {
		// Run handler in a separate goroutine so we don't block the NATS read loop
		go func() {
			if err := handler(context.Background(), m.Data); err != nil {
				log.Printf("[NATS] Error handling topic %s: %v", topic, err)
			}
		}()
	})
	return err
}

// QueueSubscribe registers a competing consumer in a Core NATS queue group.
// Core NATS remains intentionally ephemeral; OVA's SQLite reconciler owns
// durable acceptance and republishes work that has not been claimed.
func (b *Bus) QueueSubscribe(topic string, queue string, handler EventHandler) error {
	if queue == "" {
		return fmt.Errorf("queue group is required")
	}
	_, err := b.nc.QueueSubscribe(topic, queue, func(m *nats.Msg) {
		go func() {
			if err := handler(context.Background(), m.Data); err != nil {
				log.Printf("[NATS] Error handling queue topic %s: %v", topic, err)
			}
		}()
	})
	if err != nil {
		return err
	}
	return b.nc.Flush()
}

func (b *Bus) Connected() bool {
	return b != nil && b.nc != nil && b.nc.IsConnected()
}

func (b *Bus) Flush(ctx context.Context) error {
	if b == nil || b.nc == nil {
		return fmt.Errorf("NATS connection is unavailable")
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(5 * time.Second)
	}
	return b.nc.FlushTimeout(time.Until(deadline))
}

func (b *Bus) Close() {
	if b.nc != nil {
		b.nc.Close()
	}
}
