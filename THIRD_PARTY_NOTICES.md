# License map and third-party notices

This is a mixed-license review package. The root evaluation draft does not override the notices below or other applicable rights.

| Material | Applicable notice |
| --- | --- |
| `samples/distributed-queue/` | MIT; original OVA notice retained in that directory |
| Eligible original OM extracts and portfolio text | Root draft evaluation terms; finalize before publication |
| Screenshot | OM interface with synthetic data; product marks remain reserved; depicted third-party UI components retain their own licenses |
| External Go dependencies | Their upstream licenses, documented below and supplied in `licenses/` |

The Go samples do not vendor dependency source. The queue sample downloads public modules identified by its `go.mod` and integrity checks in `go.sum`. A generated dependency-license inventory is included in `licenses/dependencies.md`; original upstream license files are copied alongside it. The workspace-protection and project-loading samples have no third-party package dependencies.

The screenshot was rendered from OM's existing UI test fixtures. React, Lucide, and other UI dependencies used to render it are not included as source in this package. Their relevant notices are included in `licenses/`. Product names in the documentation identify systems and tools, not endorsements.

The source repositories' historical licensing and distribution scope require review before publication. This folder does not purport to revoke MIT/Apache permissions already granted elsewhere.
