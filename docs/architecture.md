# Architecture and boundaries

![Build workflow and distributed execution](../media/build-flow.png)

The diagram shows two related paths. It is an architecture illustration, not a statement about the size or availability of a particular deployment.

## OpenMission: application delivery

The user reviews a story and acceptance criteria before approving a build. The controller coordinates acceptance-test generation, implementation, and validation in project-specific workspaces. Validation evidence, rather than an agent's claim alone, determines reported outcomes. Test quality is still a limitation: a passing test proves only what it checks.

Local OM uses a Go API, a React/TypeScript UI, SQLite project state, and Docker-managed build/test runtimes. Hermes supplies the agent session and tools, with configured model endpoints providing inference. Security reporting uses Trivy and is informational in the described workflow.

Tasks are a separate path for changing an existing application. They keep their own run history and include regression validation. They do not redefine the original story-build flow.

## OVA: distributed cognitive execution

In the cloud-oriented path, an OM caller submits work to OVA. OVA persists the execution before publishing a dispatch identifier through Core NATS. The implemented lane model includes Thinking and Coding, with queue groups for competing workers. Owner-specific subjects support session affinity.

Workers must claim eligible work and hold the relevant lease before execution. Thinking work can call inference directly; Coding work uses a Hermes session. Metered inference is attributed to the active execution. Results and usage are stored and retrieved through the API.

Core NATS messages are ephemeral. Durable acceptance, claims, leases, and reconciliation belong to OVA's service and SQLite state, not the broker adapter in this portfolio. Reconciliation can republish work; duplicate dispatch is possible. No exactly-once guarantee for external side effects is claimed.

Provider abstractions and illustrative multi-domain/multi-agent extensions exist in the source material, but this package does not claim that all proposed agent/domain combinations are deployed or validated.

## What this folder implements

| Area | Included | Deliberately outside this package |
| --- | --- | --- |
| Messaging | NATS adapter and loopback competing-consumer test | Durable claims, leases, metering, reconciliation |
| Workspace protection | Protected-path snapshot/restore mechanics | Full worker lifecycle, upload-manifest validation, process sandboxing |
| UI loading | Latest-request selection and cancellation | Full React application, API, background polling |
| Persistence | Case study and recorded synthetic results | Complete persistence implementation and project database |

`openmission-cloud` contains the cloud-oriented application integration. `openmission-deploy` contains packaging/deployment materials. Neither source tree is needed to run these samples.

[Editable diagram source](../media/build-flow.svg) · [Diagram rendering script](../media/render-diagram.mjs)
