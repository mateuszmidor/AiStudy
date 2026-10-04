## Context

This change introduces a new CLI tool (`hwinfo`) for reading and displaying CPU information on Linux (amd64). The tool targets a single-shot execution model: it reads CPU data via the `gopsutil/v4/cpu` library and prints formatted output to stdout. There is no existing codebase — this is a greenfield project at the repository root.

## Goals / Non-Goals

**Goals:**
- Create a single-file Go CLI (`main.go`) that prints 7 fixed CPU fields
- Use `github.com/shirou/gopsutil/v4/cpu` as the sole data source
- Output human-readable plain text with a header, fixed labels, and strict field ordering
- Handle library errors gracefully (critical failure → exit 1, partial data → show `N/A`)
- Target Go 1.26+ and Linux amd64

**Non-Goals:**
- No CLI flags, arguments, or configuration
- No interactive mode or daemon behavior
- No cross-platform support (Linux amd64 only)
- No output formats beyond plain text (no JSON, YAML, etc.)
- No unit tests in this change (single-file scope)

## Decisions

1. **Single-file layout** over multi-file package: Keeps the project minimal — the total logic is ~100 lines. No benefit from splitting into packages.
2. **gopsutil/v4/cpu** over raw procfs parsing: Provides a stable Go API that abstracts over `/proc/cpuinfo`. Avoids brittle string parsing.
3. **Averaged CPU speed** over per-core display: gopsutil returns per-CPU info. Averaging across all logical cores gives a single representative value, matching the spec.
4. **`fmt.Printf` over structured logging**: Simplest approach for a single-shot CLI. No need for a logging framework.
5. **No rounding on MHz values**: Raw library values printed as-is per spec requirement.
6. **Exit 1 on critical error with no partial output**: If gopsutil fails entirely, print nothing to stdout and exit non-zero. Avoids misleading partial display.

## Risks / Trade-offs

- **[Risk] gopsutil v4 API stability**: The v4 module is relatively new. If the API changes, `go.mod` pins the version. Low risk.
- **[Risk] N/A for individual fields vs. complete failure**: If one field is empty but the library returns successfully, we show `N/A` for that field. If the library call itself fails, we exit 1 with no output. The line is clear.
- **[Trade-off] No tests**: Skipping tests keeps this change minimal. Tests should be added in a follow-up.

## Migration Plan

No migration needed — this is a new tool. Users clone the repo and run `go run` or `go build`.

## Open Questions

None — all design decisions are settled by the spec.
