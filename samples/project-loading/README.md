# On-demand project loading

A framework-independent TypeScript helper extracted from OM's UI. Governed by the root [draft evaluation terms](../../LICENSE).

## Run

```sh
npm test
```

Requires Node.js 22.18+; no npm dependencies or install step are needed. `package-lock.json` records the dependency-free package. Node may print an experimental mock-timer warning; it does not indicate test failure.

## Read

- [`src/onDemand.ts`](src/onDemand.ts): one in-flight resource request, cancellation, deduplication, and stale-result suppression.
- [`tests/onDemand.test.ts`](tests/onDemand.test.ts): three behavior tests for on-demand reads, rapid selection, canceled errors, and retry.

Only concurrent requests for the same key share an operation. Completed results are not cached. Switching keys aborts the previous request and prevents its result from reaching the caller even if transport cancellation is ignored. The helper does not define polling, UI presentation, or persistence.

The implementation is copied without behavior changes. One application-integration test that inspected private React source files is omitted, along with its filesystem import. The three standalone behavior tests are retained.
