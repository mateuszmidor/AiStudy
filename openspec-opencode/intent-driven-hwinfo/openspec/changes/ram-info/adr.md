# ADR Review Manifest

- Status: completed
- Review date: 2026-10-13

## Review Summary

ADR review completed for change `ram-info`. `<repo>/adr/` did not exist before this change, so there were no in-force ADRs to walk, no supersession graph, and no prior ADR was edited. One decision in `design.md` met the bar for a durable, repository-level record: the commitment to a memory data source, which outlives this change, affects later memory-details work, and is not captured by any ADR because none existed. One new repository-level ADR was created. The tool-level output contract (single labeled line, 1024³ **GB** with truncation, in-band `N/A` at exit 0) was considered and intentionally left out of the ADR layer, on the decision that it belongs to the spec and the glossary; this is recorded here for the audit trail rather than treated as an oversight.

## In-Force ADRs Reviewed

- None — `<repo>/adr/` did not exist when this review began, so there were no in-force repository-level ADRs constraining the design and nothing to supersede.

## New Durable ADRs Created

- `adr/0001-use-gopsutil-for-total-physical-memory.md` — the standing commitment to `gopsutil/v4/mem` as the memory data source, pinned at `v4.26.9`, with the `x/sys` syscall and the `sysctl` subprocess as the rejected alternatives.
