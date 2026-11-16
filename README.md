# Charlie Kuharski — Engineering Portfolio

**Technology leadership, cloud platforms, and AI application delivery.**

I am the sole creator and contributor to **OpenMission (OM)** and **OpenVirtualAgent (OVA)**. This portfolio presents selected engineering work alongside the product and architectural decisions behind it. It is a review package, not a distribution of either complete product.

My professional background includes technology leadership in financial services and leading approximately 20 people, including managers, across North America, England, and India. OM/OVA are my independent projects; they are not presented as work performed by those teams.

## The products

**OpenMission** coordinates application development from a user request through a reviewed story, acceptance criteria, generated tests, implementation, and validation. Its workspace brings build progress, source, test evidence, application controls, and security findings together.

**OpenVirtualAgent** supplies queued cognitive execution for the cloud-oriented path: durable execution records, distributed dispatch, worker claims and leases, provider execution, and inference metering. Local OM also has a direct Hermes path; OVA is not required for every local build.

I use AI-assisted development tools. Product direction, architecture, integration, debugging decisions, and acceptance of changes remain my responsibility.

![OpenMission build progress rendered with synthetic project data](media/overview.png)

*Actual OM interface, rendered against local synthetic API fixtures. Names, progress, and findings are demonstrations, not customer data or a live production run.*

## Five-minute review

1. **One minute:** inspect the [architecture](docs/architecture.md) and [flow diagram](media/build-flow.png).
2. **Two minutes:** choose a code sample below and read its implementation and tests.
3. **One minute:** review the [engineering decisions](docs/engineering-decisions.md), including limitations.
4. **One minute:** read the [persistence case study](docs/build-persistence-case-study.md) or [technical leadership](docs/technical-leadership.md).

| Sample | Evidence | Start here |
| --- | --- | --- |
| Distributed queue | Competing consumers using an embedded NATS server | [Code, test, and scope](samples/distributed-queue/README.md) |
| Workspace protection | Snapshot/restore behavior, artifact exclusions, file permissions, failure handling | [Code, tests, and scope](samples/workspace-protection/README.md) |
| Project loading | In-flight deduplication, cancellation, stale-response suppression, retry | [Code, tests, and scope](samples/project-loading/README.md) |

## Demos and posts

More samples and demonstrations are coming. Follow the work here:

- [YouTube channel — videos and demonstrations](https://www.youtube.com/channel/UCYSKROYZ2riNfj9zyVbFB3Q/videos)
- [LinkedIn post](https://www.linkedin.com/feed/update/urn:li:activity:7526869391683141632/)

## Run the samples

Requirements: Go 1.25 or newer, and Node.js 22.18 or newer. Go race checks require a supported native toolchain. macOS is the verified environment for the filesystem sample's permission behavior. Windows is not validated here; the Unix umask test is excluded there.

From this folder:

```sh
(cd samples/distributed-queue && go test -race -count=1 ./...)
(cd samples/workspace-protection && go test -race -count=1 ./...)
(cd samples/project-loading && npm test)
```

The first Go command may download public dependencies. The queue test launches its own broker bound to loopback; no installed NATS service is required. The loader sample has no npm dependencies. No Docker, model server, credentials, private repository, or OM installation is needed.

See [verification](docs/verification.md) for the recorded checks, and [provenance](docs/provenance.md) for source revisions and extraction changes.

## Further discussion

A product demonstration, architecture walkthrough, and additional private source review are available upon request. These examples deliberately omit complete orchestration logic, agent prompts, production configuration, and operational data.

## Licensing

This package has **mixed licensing**, not a blanket open-source license. The OVA queue extract retains its [MIT license](samples/distributed-queue/LICENSE). The remaining eligible original materials have [draft hiring-evaluation terms](LICENSE); third-party rights are described in [the notices](THIRD_PARTY_NOTICES.md). Existing rights granted under another applicable license are not withdrawn.

Owner preparation: complete the [pre-publication checklist](docs/pre-publication-checklist.md) before uploading this folder.
