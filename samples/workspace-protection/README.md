# Workspace protection

A focused extract of OM's protected-path snapshot and restoration mechanism. Governed by the root [draft evaluation terms](../../LICENSE).

## Run

```sh
go test -race -count=1 ./...
```

Requires Go 1.25+. There are no third-party Go dependencies; consequently no `go.sum` is needed. Tests use temporary directories and synthetic content. File-mode/umask behavior is tested on Unix; Windows has not been validated.

## Read

- [`protection.go`](protection.go): capture protected entries, detect changes, restore content/permissions/symlinks, then verify restoration.
- [`protection_test.go`](protection_test.go): generated reports remain writable; restoration failure is returned.
- [`protection_unix_test.go`](protection_unix_test.go): regular-file content and permissions survive restoration under a restrictive umask.
- [`fixtures_test.go`](fixtures_test.go) and [`paths.go`](paths.go): minimal standalone support.

The classification protects controller metadata, non-allowed tests, and selected sensitive-file names while ignoring generated test-result artifacts. This is not a hostile-process sandbox. It assumes coordinated workspace access and cannot undo network activity, prevent secret reads, or guarantee safety against concurrent filesystem replacement attacks. Snapshot memory grows with protected content.

## Extraction changes

The package is renamed from `server` to `protection`. The original upload-manifest validator is replaced by an explicit `protectUploads` parameter supplied by the caller; the retained tests pass `false`. The sample therefore does not demonstrate upload-manifest validation. Test fixtures no longer change OM's global workspace root.

Two unused reporting/checkpoint comparison helpers are omitted. Three direct snapshot/restore tests are retained; full acceptance-worker integration tests requiring private controllers are omitted. The Unix-only test is separated with a build constraint. No complete agent prompt, receipt protocol, or controller implementation is included.
