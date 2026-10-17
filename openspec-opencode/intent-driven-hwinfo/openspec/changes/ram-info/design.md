## Context

This is a greenfield repository: the only tracked artifact is `go.mod` (`module hwinfo`, `go 1.27`) and there is no `go.sum`, no source, no `README`, and no `adr/` directory. Because `adr/` does not exist, there are **no in-force ADRs** constraining this design and nothing to supersede; this change's `adr` step will create the first one. Toolchain in place: `go1.27.1 darwin/arm64`.

The target machine reports `hw.memsize` = 68719476736 bytes, exactly 64 × 1024³, and About This Mac shows "64 GB"; dividing by 1000³ would print 68.71. That is why the display convention is 1024³ labeled `GB`, and why the truncation consequence is invisible here.

On darwin, `mem.VirtualMemory()` reaches `.Total` through `getHwMemsize()`, which is `unix.SysctlUint64("hw.memsize")`. It gets there only after `common.NewSystemLib()` (a purego `dlopen` of libSystem), a `HostStatistics(HOST_VM_INFO)` call, and a `Dlsym("vm_kernel_page_size")` — so the mach VM plumbing runs alongside the read but does not compute the number. `.Total` is the sysctl value.

Stakeholders: the repository owner, who is the only person who commits. The Linux twin `../hw` is a reference implementation, not a dependency.

### System shape (lightweight C4)

```
   user runs hwinfo
        |
        v
  +----------------+     one line     +----------------------+
  |    main.go     | <--------------- |   internal/ram       |
  |  orchestrate   |                  |  FormatTotalLine()   |
  |  + print       | ---- call -----> |  Total()             |
  +----------------+                  +----------------------+
        |                                        |
        |                                bytes / error
        v                                        v
  stdout, exactly one line          +------------------------+
  Total Memory: 64 GB               |  gopsutil/v4/mem       |
  or Total Memory: N/A              |  VirtualMemory()       |
  (always exit status 0)            |   -> hw.memsize        |
                                    +------------------------+
                                                 |
                                                 v
                                    macOS: libSystem via purego
                                    dlopen; sysctl hw.memsize
```

- **Boundary**: one runnable container — a single binary. There is no service, no daemon, no config, no persistence.
- **Responsibility split**: `main.go` owns process concerns (when to print, what goes to stdout, exit status); `internal/ram` owns the hardware read and the bytes-to-line conversion; `gopsutil/v4/mem` owns talking to macOS.
- **Key relationship**: the only external dependency in the picture is the library boundary, and it is one call deep.
- **Assumption**: the failure modes worth handling all surface as a returned error, because the syscall underneath is what supplies the number.
- **Open questions**: none remain at this level — the diagram is deliberately one level deep because a component view of a three-file package would have nothing to separate.

## Goals / Non-Goals

**Goals:**
- Implement the `ram-info` capability exactly as specified: one line on stdout, `N/A` for an error or a zero total, silent stderr, exit 0, arguments ignored.
- Keep the hardware read behind one seam so the conversion arithmetic can be tested without touching the OS.
- Make the `gopsutil/v4` commitment concrete enough for the `adr` step to record it durably.

**Non-Goals:**
- Memory details (`Available`, `Used`, `Free`, `Wired`, `UsedPercent`), CPU reporting, a header block, unit flags, humanized formatting, configuration, or JSON output.
- Claiming or verifying Linux. The code compiles there and would read `/proc/meminfo`, but no behaviour is promised.
- An acceptance-test harness: this workspace ships Python and JS runners, not Go ones, so verification is a unit test plus one manual run.
- Continuous integration.

## Decisions

1. **Data source: `gopsutil/v4/mem`, via `mem.VirtualMemory().Total`.** Alternatives considered: `golang.org/x/sys/unix.SysctlUint64("hw.memsize")` — the fewest moving parts, one direct dependency, no purego, no dlopen, ~5 lines, and already in the module cache — rejected because it returns only the total, so a later memory-details slice would need a second data source, and because it diverges from the twin; and `os/exec` with `sysctl -n hw.memsize` — zero dependencies, but it spawns a process, parses stdout as text, and adds a failure mode a syscall cannot have (binary not on `PATH`). The wider stat struct is the reason to accept the heavier dependency.

2. **Module layout: three source files, with the conversion inside `internal/ram`.** `main.go` holds orchestration and printing; `internal/ram/ram.go` holds `func Total() (uint64, error)` and `func FormatTotalLine(total uint64, available bool) string`; `internal/ram/ram_test.go` tests both. This revises assumption A1, which had `main` owning the conversion — the deciding argument is that the interesting logic is the divisor and the truncation, whose real-world consequence is invisible on this machine (64 GiB divides exactly), so it must be reachable from the layout's only test file. Printing stays in `main`, so the output contract still has exactly one writer. Alternatives considered: conversion in `main` tested from a new `main_test.go` — five files, keeps `internal/ram` single-purpose, rejected to hold the requested four-file layout; and no tests at all — what the twin's first change chose, rejected because it would leave the truncation rule unverified anywhere.

3. **One formatter produces the whole line, including the `N/A` form.** `FormatTotalLine(total, available)` returns either `Total Memory: <n> GB` or `Total Memory: N/A`, so the label `Total Memory: ` has exactly one home. Alternative considered: an exported `UnavailableLine` constant plus a numeric-only formatter — rejected because the label would then exist in two places and the `N/A` form would have no test. The availability policy itself (`err != nil || total == 0`) is computed in `main`, because zero-means-unavailable is a display rule, not a read failure, and overloading the error channel with it would hide the difference.

4. **Read seam: an unexported package-level function variable, swapped by tests.** `var readTotal = func() (uint64, error) { ... }` delegating to `gopsutil`, with `Total()` calling it. `ram_test.go` replaces it to assert byte passthrough and error propagation, then restores it. Alternatives considered: `totalFrom(read func() (uint64, error))` with no mutable package state — cleaner in principle, but the test then never exercises the exported entry point, so a mistake in `Total()` itself would go unnoticed; and a reader interface — ceremony for a single call site.

5. **Nothing is ever fatal.** `main` prints one line and returns, so the process exits 0 on both the happy path and the unavailable path, and writes nothing to stderr. There is no logging framework, no `os.Exit` call, and no error text anywhere. Rationale: the spec pins all four behaviours, and an in-band signal is what keeps the tool a single-line program.

6. **No argument handling and no context plumbing.** `main` ignores `os.Args` rather than using `flag`, and calls the context-free `mem.VirtualMemory()` rather than `VirtualMemoryWithContext(ctx)`, because a one-shot syscall read has nothing to cancel and no flag is specified.

7. **Dependency install pins a version already in the local module cache.** Use `go get github.com/shirou/gopsutil/v4@v4.26.9` followed by `go mod tidy`, which creates `go.sum`. The cache already holds `v4.26.9` (also `v4.26.6`, `v4.26.4`, `v4.25.12`) and `purego@v0.11.1`, so the install does not have to reach the network. The pin is also what the ADR will record as the commitment. Alternative: `go get ...@latest`, rejected because it needs the network and silently moves a durable decision.

8. **Diagram format: one lightweight ASCII sketch, no strict C4 levels** (see Context). A component-level view of a three-file tool would separate nothing, and strict labelling would add ceremony without answering a question.

## Risks / Trade-offs

- [Risk] The darwin path uses purego to `dlopen` libSystem at runtime; if that fails, the read errors and the user sees `N/A` with exit 0, which is indistinguishable from "no memory read". -> Mitigation: this is the specified behaviour, the failure is in-band rather than fatal, and the manual verification run proves the happy path on the target machine.
- [Risk] `go get` may need the public proxy and checksum database, since `GOPROXY` and `GOSUMDB` point at `proxy.golang.org` and `sum.golang.org`. -> Mitigation: pin to the cached version; if the checksum database is unreachable, run the install step with `GOSUMDB=off` for that module only and re-enable it afterwards.
- [Risk] Truncation makes a machine one byte below a boundary print a smaller number (`17179869183` bytes prints `15 GB`, not `16`). -> Mitigation: accepted deliberately for About This Mac parity; the spec scenario and the unit test make the rule explicit and prevent drift.
- [Risk] `GB` meaning 1024³ is technically a `GiB`, so the label misleads an SI-minded reader. -> Mitigation: the glossary pins the meaning, the label was chosen so the number matches what macOS itself displays, and the ADR can record the convention.
- [Trade-off] The capability is change-scoped, so a future CPU or memory-details capability will restate the tool-level output conventions. -> Mitigation: recorded in the proposal; the ADR is the right place to state the convention once so a later split has a source.
- [Trade-off] No acceptance harness, so end-to-end correctness rests on one manual run against `sysctl -n hw.memsize` and About This Mac. -> Mitigation: the tasks list that run as an explicit, verifiable step, and the unit test covers the arithmetic that the manual run cannot exercise on this machine.

## Migration Plan

Nothing to migrate: this is a new binary with no consumers, no stored data, and no configuration. Deployment is `go build`; rollback is deleting the files.

Order of work: install the pinned dependency (`go.mod`/`go.sum`), implement `internal/ram`, implement `main.go`, run `go test ./...`, build, then run the binary and compare its output with `sysctl -n hw.memsize` and About This Mac. Expected output on this machine: `Total Memory: 64 GB`.

## Open Questions

None blocking. This design closes the three items the proposal deferred:

- Test scope: one Go unit test file (`internal/ram/ram_test.go`) covering the conversion truncation cases plus the injected read; no acceptance-test harness.
- API names: `Total() (uint64, error)` and `FormatTotalLine(uint64, bool) string` in `package ram`.
- Injection seam: the unexported `readTotal` function variable.

ADR candidates for the `adr` step to distill: the standing commitment to `gopsutil/v4/mem` as the memory data source (including the pinned version), and — optionally — the 1024³-labeled-`GB` display convention with truncation, which is tool-level and deliberately diverges from the twin's `GiB`. No in-force ADR needs revisiting, since none exists yet.
