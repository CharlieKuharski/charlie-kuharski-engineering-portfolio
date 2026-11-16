# Distributed queue adapter

Extracted from OVA's Core NATS adapter. Licensed under [MIT](LICENSE).

## Run

```sh
go test -race -count=1 ./...
```

Requires Go 1.25+ and a toolchain supporting the race detector. Public dependencies are pinned in `go.mod`/`go.sum`. The test starts an embedded broker on `127.0.0.1` with a temporary port and stops it afterward. It requires no Docker or existing NATS service.

## Read

- [`bus.go`](bus.go): publish, subscribe, competing-consumer subscription, flush, and close.
- [`handler.go`](handler.go): minimal handler contract replacing the private domain-package import.
- [`bus_test.go`](bus_test.go): 20 publications and two subscribers in one queue group, expecting 20 handler deliveries.

A queue group dispatches a publication to one available subscriber in that group. This test does not prove durable delivery or exactly-once external effects. Core NATS is ephemeral. In the complete service, SQLite claims, leases, and reconciliation are separate mechanisms and are not included here.

The adapter spawns a goroutine for each callback and has no bounded worker pool. `Publish` accepts but does not use its context; `Flush` applies a timeout/deadline. These behaviors are preserved from the reviewed source rather than silently redesigned for the portfolio.

The adaptation removes the original domain import, includes the required handler type, and retains the behavioral test. [Provenance](../../docs/provenance.md).
