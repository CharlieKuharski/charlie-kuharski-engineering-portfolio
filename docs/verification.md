# Verification record

Prepared 2026-11-16 on macOS ARM64 using Go 1.26.3 and Node.js 22.22.2.

| Sample | Command (from its directory) | Result |
| --- | --- | --- |
| distributed-queue | `go test -race -count=1 ./...` | PASS: one embedded-broker test |
| workspace-protection | `go test -race -count=1 ./...` | PASS: three filesystem tests, including Unix permissions |
| project-loading | `npm test` | PASS: three request lifecycle tests |

The Node runner emitted its experimental MockTimers warning; tests completed successfully. The Go files were formatted with gofmt. These checks cover the extracted samples, not the complete products or all operating systems.

## Visuals and package checks

- The screenshot uses the actual locally built OM interface against synthetic, intercepted API responses. All requests to other origins were blocked. No deployed API or project data was used. Rendering completed without browser page errors.
- The screenshot and architecture diagram were visually inspected at their native 1440-pixel width.
- The SVG is editable. `media/render-diagram.mjs` renders its PNG with Playwright. Supply an installed Playwright module through `PLAYWRIGHT_MODULE`, or install Playwright in a separate temporary tool environment; a compatible Chromium installation is required. This optional renderer is not a sample dependency.
- Local Markdown links, provenance source hashes, and absence of Git history, databases, environment files, dependency directories and common credential patterns were checked.
- No repository was initialized, no branches were changed, and no materials were uploaded.

The separately reported persistence case study is local source-project evidence, not a benchmark reproduced by these samples. Its deployment validation remains pending.
