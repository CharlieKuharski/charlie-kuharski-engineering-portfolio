# Case study: project loading during an active build

**Status: implemented and locally tested in an uncommitted working-tree change at preparation time; deployment verification pending.** This narrative is not a claim about released behavior. The implementation is intentionally not included in this review package.

## Symptom and investigation

Project loading was responsive when builds were idle but slower while another project built. Local HTTP observations returned compact project responses in roughly 3.7–3.9 seconds during the investigation. Those observations were not a controlled benchmark.

The selected-project API was already compact. Inspection showed that each recorded build event appended in memory, then invoked the general project save path. That path deleted and rewrote the affected project's stored records, including accumulated events, while holding the shared state lock. Overview reads needed the same lock. The persistence work therefore coupled an unrelated reader to a busy writer.

## Implemented change

- Batch routine build/task progress every five minutes; keep the UI live from memory.
- Insert events after the last committed sequence, and update metadata rather than replacing accumulated history.
- Save completion, failure, pause/cancel, and user changes immediately.
- Separate read-state synchronization from mutation serialization, permitting consistent reads during database I/O while ordering writes.
- Extract task events from embedded task JSON through a transactional migration.
- Retain failed batches for retry and flush pending work during graceful shutdown.

The selected tradeoff permits up to five minutes of unflushed routine progress to be lost in an unexpected crash under normal storage operation. Database failures can extend that window. This choice does not make delayed writes appropriate for every application; it matches this application's agreed durability requirements.

## Evidence

A synthetic in-process test opened project B's Overview 100 times while project A had 10,000 existing events and recorded additional activity. The recorded mean was **180.83 microseconds idle and 208.563 microseconds while building** on the local test machine. These numbers exclude HTTP transport and browser rendering and cannot be compared directly with the earlier 3.7–3.9-second HTTP observations.

A separate failure-injection test held a database write open and verified that Overview and sidebar reads completed before the write was released. Other tests covered batching with simulated time, append-only event persistence, failed-commit retries, ordered deletion, migration rollback, and final worker-event flushing on shutdown.

The source project's full Go suite and focused race checks passed. UI checks included 150 unit tests and browser checks at desktop and two mobile widths. An existing unreachable-code warning remained in an unrelated product-UI planning file. These are source-project results, separate from the standalone sample checks in [verification](verification.md).

## Remaining validation

Repeat latency measurements against the deployed API during an actual build, including tail latency and multiple projects. Measure database transaction time and pending-event volume under prolonged load. Do not treat the synthetic figures as a production SLA or a claim of unlimited concurrency.
