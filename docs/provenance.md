# Source provenance

Prepared 2026-11-16. Only the listed extracts are included. Repository history, private controllers, operational data, and source-project credentials are not part of this folder.

| Repository | Source | Revision | Local status |
| --- | --- | --- | --- |
| OpenVirtualAgent | `pkg/bus/nats/bus.go` | `5a1ed0070604` | Matches recorded commit |
| OpenVirtualAgent | `pkg/bus/nats/bus_test.go` | `5a1ed0070604` | Matches recorded commit |
| OpenVirtualAgent | `LICENSE` | `5a1ed0070604` | Matches recorded commit |
| openmission | `internal/server/worker_workspace_protection.go` | `dbf509b08b0d` | Matches recorded commit |
| openmission | `internal/server/worker_workspace_protection_test.go` | `dbf509b08b0d` | Matches recorded commit |
| openmission | `ui/src/onDemand.ts` | `dbf509b08b0d` | Matches recorded commit |
| openmission | `ui/tests/onDemand.test.ts` | `dbf509b08b0d` | Matches recorded commit |

Exact source hashes and destination mappings are in [provenance.json](provenance.json). Extracted-source hashes describe original files before documented adaptations; they are not hashes of the adapted destination files.

## Adaptations

- Queue: remove private domain-package import; supply the compatible `EventHandler` signature locally; retain original queue test and MIT notice.
- Workspace: rename package, accept an explicit upload-protection flag, omit two unused reporting/checkpoint helpers, provide temporary-directory test helpers, retain three direct behavior tests and split the Unix umask test. Full worker integration tests are not included.
- Loader: implementation unchanged; retain three standalone tests and omit the source-inspection integration test.
- Packaging: standalone module/package manifests, sample READMEs, diagram, and review documents are new.

## Case-study provenance

The build-persistence implementation was an uncommitted OM working-tree change when this folder was prepared. The case study cites local synthetic results from that work; it is not attributed to the committed source revision above. The underlying implementation and test logs are deliberately not copied.

## Visual provenance

The screenshot is a new local rendering of the OM UI against synthetic API data. It contains no live project response, customer record, or private conversation. The SVG diagram is newly authored from the reviewed architecture documentation; it is illustrative rather than a deployment inventory.
