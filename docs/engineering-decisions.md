# Engineering decisions

## 1. Separate dispatch from durable execution

**Problem:** multiple workers need to receive available work without every worker executing each publication.

**Decision:** use Core NATS queue groups for dispatch and keep durable execution state in OVA. The extracted adapter is deliberately small. Its queue test starts an embedded broker and verifies that 20 publications produce 20 handler deliveries across two subscribers in the same group.

**Tradeoff:** simple dispatch does not provide durable replay. The complete service needs claim checks, lease handling, and reconciliation. The adapter alone does not implement backpressure, bounded handler concurrency, durable acknowledgments, or exactly-once processing. Its ordinary `Publish` method does not honor context cancellation; `Flush` uses a deadline. These are limits of the extract, not guarantees hidden behind the interface.

[Read the sample](../samples/distributed-queue/README.md).

## 2. Protect controller-owned files while allowing useful work

**Problem:** a coding worker needs to modify application files without silently rewriting the specification or other protected evidence. Generated diagnostic files must remain writable.

**Decision:** snapshot selected protected paths, allow explicitly assigned test paths and mutable paths, then restore and verify protected content. Exclude reports, traces, screenshots, and similar result artifacts. Treat restoration failure as an error.

**Tradeoff:** snapshots consume memory proportional to the protected content and require coordination with the worker. This is a verification/restoration mechanism, not a security sandbox: it cannot prevent exfiltration, undo arbitrary external side effects, or eliminate concurrent filesystem attacks. Isolation and exclusive workspace ownership remain separate concerns.

The extract accepts an explicit upload-protection flag instead of importing OM's complete manifest validation. That adaptation is documented and tested only within the sample's stated scope.

[Read the sample](../samples/workspace-protection/README.md).

## 3. Do not let late responses change the selected project

**Problem:** a user switches from project A to B while A's network request is still pending. A can finish last and overwrite B's screen.

**Decision:** cancel the previous resource request; also suppress its result if the transport ignores cancellation. Deduplicate requests for the same in-flight key. Explicit requests after completion fetch again rather than returning a hidden cache.

**Tradeoff:** there is no completed-result cache, polling policy, or error presentation in this helper. Those belong to its caller. The extracted tests cover out-of-order completion, cancellation, errors, and retry without relying on a live backend.

[Read the sample](../samples/project-loading/README.md).

## 4. Diagnose the backend before changing the UI

A later project-loading delay occurred while a build was active. The compact read API already fetched the selected workspace. Investigation found database persistence rewriting accumulated event history while holding a shared state lock. That required a persistence fix rather than another loading animation or client cache.

The implementation and evidence are described in the [case study](build-persistence-case-study.md). Deployment verification remains pending; no production improvement percentage is asserted.
