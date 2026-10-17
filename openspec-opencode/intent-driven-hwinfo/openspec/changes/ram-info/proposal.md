## Why

macOS has no parse-free command that reports installed **total physical memory** the way the system itself describes it: `sysctl -n hw.memsize` prints raw bytes and `system_profiler SPHardwareDataType` prints a verbose report, so reproducing the "64 GB" that About This Mac shows means knowing macOS's 1024³ convention. **hwinfo** begins with this slice — one labeled line — and later CPU and memory-detail reporting will join it.

## What Changes

- Add a macOS CLI that prints a single line, `Total Memory: <n> GB`, where `<n>` is **total physical memory** expressed in **GB** — bytes divided by 1024³, truncated to an integer.
- Read the total with `gopsutil/v4/mem` (`mem.VirtualMemory().Total`): `internal/ram` returns raw bytes plus an error, and `main` owns the conversion and the printing.
- On a read failure, print `Total Memory: N/A` on stdout and exit 0, keeping the output shape stable instead of signalling an error.
- Add the app skeleton: `main.go`, `internal/ram/ram.go`, `internal/ram/ram_test.go` (import path `hwinfo/internal/ram`).
- Add `github.com/shirou/gopsutil/v4` as the module's first direct dependency.
- Target macOS **by intent, not by compiler constraint** — the same code compiles on Linux and reads `/proc/meminfo`, but this change claims and verifies macOS only.
- Establish the tool-level output conventions inside this capability: one labeled line as the whole stdout contract, 1024³ **GB**, and `N/A` + exit 0 on read failure.
- Explicitly **not** in this change: memory details (`Available`, `Used`, `Free`, `Wired`, `UsedPercent`), CPU reporting, a header block, unit flags (no `--gib` or `--bytes`), humanized formatting, and argument parsing — arguments and flags are ignored.

## Capabilities

### New Capabilities

- `ram-info`: Report installed **total physical memory** as one labeled line whose number matches what macOS shows for the same machine, and report `N/A` at exit 0 when the value cannot be read.

The name is change-scoped by decision: a later CPU slice or a memory-details slice adds its own capability rather than widening this one, matching how the Linux twin (`../hw`) actually grew (`cpu-info-display`, then `cpu-temp`).

### Modified Capabilities

None. This is the first capability in this repository and `openspec/specs/` is empty.

## Impact

- **New code**: `main.go` (format and print), `internal/ram/ram.go` (read total bytes), `internal/ram/ram_test.go`. No existing code is modified — the repository currently holds only `go.mod`.
- **New dependency**: `github.com/shirou/gopsutil/v4` (the `mem` sub-package) as the module's first direct dependency. It pulls transitive modules, including `ebitengine/purego` (used to `dlopen` libSystem), `golang.org/x/sys`, and platform/test helpers. This is a durable commitment, so it is an ADR candidate to be recorded at `adr/NNNN-*.md` after design.
- **Output conventions**: the single labeled line, 1024³ **GB**, and `N/A` + exit 0 live inside the `ram-info` capability for now. If later metric capabilities duplicate them and that duplication starts to hurt, a tool-level capability can be split out in a later change.
- **Deliberate divergence from the Linux twin**: shape and unit differ — `Total Memory: 64 GB` here versus the twin's `HWInfo - Memory Information` header block and GiB — while the failure path (`N/A` + exit 0) matches verbatim. The two tools are not expected to diff clean, and that is intentional. The twin's own `add-ram-reporting` change is open and never applied.
- **Unresolved, deferred to design**: the `ram` package API surface (a `func Total() (uint64, error)` is assumed, not confirmed) and whether `Total()` calls `gopsutil` directly or takes an injected reader so the conversion is testable. Test scope is likewise design-time: one Go unit test for the bytes-to-line conversion, as this workspace ships no Go acceptance-test harness.
- **Unaffected**: no API, no configuration, no persisted data, no migration, and no breaking changes — the repository has no consumers yet.
