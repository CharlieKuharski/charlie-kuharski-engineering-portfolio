package nats

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	natsserver "github.com/nats-io/nats-server/v2/server"
)

func TestQueueSubscribeDeliversEachMessageToOneConsumer(t *testing.T) {
	server, err := natsserver.NewServer(&natsserver.Options{Host: "127.0.0.1", Port: -1})
	if err != nil {
		t.Fatal(err)
	}
	go server.Start()
	if !server.ReadyForConnections(5 * time.Second) {
		t.Fatal("embedded NATS did not start")
	}
	t.Cleanup(server.Shutdown)

	bus, err := NewBus(fmt.Sprintf("nats://%s", server.Addr().String()))
	if err != nil {
		t.Fatal(err)
	}
	defer bus.Close()
	var deliveries atomic.Int64
	handler := func(_ context.Context, _ []byte) error {
		deliveries.Add(1)
		return nil
	}
	if err := bus.QueueSubscribe("ova.execution.thinking", "ova-thinking", handler); err != nil {
		t.Fatal(err)
	}
	if err := bus.QueueSubscribe("ova.execution.thinking", "ova-thinking", handler); err != nil {
		t.Fatal(err)
	}
	for index := 0; index < 20; index++ {
		if err := bus.Publish(context.Background(), "ova.execution.thinking", map[string]int{"id": index}); err != nil {
			t.Fatal(err)
		}
	}
	if err := bus.Flush(context.Background()); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(2 * time.Second)
	for deliveries.Load() < 20 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if got := deliveries.Load(); got != 20 {
		t.Fatalf("deliveries=%d, want exactly 20", got)
	}
}
