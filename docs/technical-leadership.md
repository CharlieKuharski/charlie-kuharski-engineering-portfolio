# Technical leadership and hands-on ownership

## Independent product work

I conceived and built OpenMission and OpenVirtualAgent as the sole contributor. The products reflect decisions across workflow design, UI, backend state management, agent integration, generated-code validation, security reporting, and deployment boundaries. AI-assisted development is part of my workflow; I remain responsible for integration, evaluation, and deciding what is ready to ship.

The samples here are intentionally narrow. They let a reviewer inspect concrete behavior and tests while the full product implementation remains private. A walkthrough can connect each sample to its role in the larger system.

## Professional leadership context

My supplied professional history includes a Director role in financial services, leading approximately 20 people including managers across North America, England, and India. My responsibilities have included cloud-native application and data-science delivery, AWS automation, CI/CD, and transformation within a matrixed organization. Earlier work included a manager/architect role supporting Cargill through Daugherty Business Solutions.

These are self-reported career details. They do not imply that my employers sponsored OM/OVA, that those teams contributed to these products, or that this portfolio contains employer-owned code or operational data.

## How I approach engineering decisions

- **Make the workflow visible:** users need build steps and outcomes, not just raw agent logs.
- **Separate claims from evidence:** a model's completion claim is not equivalent to acceptance validation.
- **Prioritize the actual constraint:** the persistence case study addresses backend contention rather than masking it with UI caching.
- **Ship bounded changes:** local fixtures and explicit deployment approval help distinguish a tested implementation from a deployed outcome.
- **State limits clearly:** passing tests, a queue adapter, and a security scan each establish something specific; none establishes complete product correctness or compliance.

## Useful interview discussion

How should durability requirements change as the product scales? Where should worker ownership and cancellation live? How should generated tests be evaluated when the model writes both the implementation and tests? What should be enforced by deterministic code, and what should remain model judgment?

I am interested in roles combining technology leadership, product ownership, architecture, and substantial hands-on engineering. Product demonstration and additional private source review are available upon request.
