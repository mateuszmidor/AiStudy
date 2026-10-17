# Decisions — hwinfo RAM reporting (macOS)

Status: pre-proposal. Captured from an `/opsx-explore` session. **No code written, no
OpenSpec change created.** This file is the input for `/opsx-propose`, not an artifact.

Target: a very small macOS CLI that prints the installed RAM total, e.g. `Total Memory: 64 GB`.

## Confirmed decisions

| # | Decision | Chosen |
|---|----------|--------|
| D1 | Scope | First slice of a growing `hwinfo` tool, not a one-trick binary. Capability named for the tool so CPU and memory-detail slices can join it later, the way `../hw` grew. |
| D2 | Data source | `gopsutil/v4/mem` — `mem.VirtualMemory().Total`. |
| D3 | stdout contract | One labeled line: `Total Memory: 64 GB`. Nothing else on stdout. |
| D4 | Number convention | bytes / 1024^3, integer truncation, label `GB`. On this Mac: `68719476736` -> `Total Memory: 64 GB`, matching About This Mac. |
| D5 | Read failure | `Total Memory: N/A` on stdout, exit 0. |

### Consequences the user accepted knowingly

- **D2** is a durable dependency commitment, so this change has a real repo-level ADR
  candidate (`adr/NNNN-...md`), which is where `../hw`'s ADR-0001/ADR-0002 live.
- **D2** also matches `../hw`'s Linux approach, so the twins read memory identically.
- **D2** returns `Available`, `Used`, `Free`, `Wired`, `UsedPercent` alongside `Total`,
  so a later "memory details" slice is a wider print, not a new data source.
- **D4** is currently invisible: `hw.memsize` is exactly `64 x 1024^3` here, so truncation
  and the divisor choice only show up on a machine whose byte count divides unevenly.

## Assumptions stated, not challenged

- **A1** `internal/ram` returns **raw bytes + error**; `main` owns formatting and printing.
  Derived from the phrasing "the memory reading function under `./internal/ram/`".
- **A2** *(corrected during the session)* The tool is macOS-only **by intent**, not by the
  compiler. With `gopsutil/v4/mem` the same code compiles on Linux and reads `/proc/meminfo`,
  which is why `../hw` can share the approach. An earlier note that darwin-only was a compile
  constraint applied to the rejected `x/sys` option, not to the chosen one.
- **A3** The intent-driven workflow is the real subject; the RAM app is the vehicle.

## Requested file layout

```
intent-driven-hwinfo/
    go.mod                      module hwinfo, go 1.27        (already present)
    main.go                     package main                  -- format + print
    internal/ram/ram.go         package ram                   -- read total bytes
    internal/ram/ram_test.go    package ram                   -- the seam tests
```

Import path is `hwinfo/internal/ram`. `internal/` is used to name the seam (one function
touching the OS behind one line of `main`), not to hide the package from external consumers.

## Grounding evidence

- `go.mod`: `module hwinfo`, `go 1.27`. Toolchain present: `go1.27.1 darwin/arm64`.
- `openspec/config.yaml` context: "Domain: macos cli tooling"; use `senior-dev` for all src
  work and `senior-qa` for acceptance tests; do not batch their tasks.
- `openspec/config.yaml` rules: proposal **must** use `grill-me`, `adversarial-authoring`,
  `glossary`; design **must** use `c4-diagrams` + `glossary`; adr **must** use
  `architectural-decision-records`. `spec-as-source` is commented out. So the proposal is the
  most expensive artifact in this change, not the code.
- `openspec list` / `list --specs`: no changes and no specs — clean slate.
- This machine: `sysctl -n hw.memsize` -> `68719476736` (= 64 GiB exactly);
  `/1000^3` would print `68.71`. macOS 27.0.1.
- `gopsutil@v4.26.9/mem/mem_darwin.go`: `getHwMemsize()` at line 17 is
  `unix.SysctlUint64("hw.memsize")`; `VirtualMemoryWithContext()` reaches it at line 118
  *after* `common.NewSystemLib()` (purego `dlopen` of libSystem), `HostStatistics(HOST_VM_INFO)`
  and `Dlsym("vm_kernel_page_size")`. So `.Total` is the syscall, reached via mach VM plumbing
  that does not touch the requested number.
- `gopsutil@v4.26.9/go.mod` requires `ebitengine/purego v0.11.1`, `go-cmp`, `plan9stats`,
  `perfstat`, `testify`, `go-sysconf`, `wmi`, `x/sys`; indirect `go-ole`, `numcpus`, `yaml`.
- `golang.org/x/sys@v0.49.0` is in the module cache; `unix.SysctlUint64` lives at
  `unix/syscall_bsd.go:454` and hardcodes an 8-byte read.
- `../hw` (Linux twin): CPU reporting done, 2 ADRs, several archived changes, and an open
  never-applied change `add-ram-reporting` whose design chose `gopsutil/v4/mem`, `GiB`,
  a `HWInfo - Memory Information` header block, and `N/A` + exit 0.

## Rejected options, with reason

- **`golang.org/x/sys/unix` — `unix.SysctlUint64("hw.memsize")`.** Fewest moving parts: one
  direct dep, no cgo, no subprocess, ~5 lines, already cached. Recommended on those grounds;
  declined in favour of `../hw` parity and the wider stat struct.
- **`os/exec` + `sysctl -n hw.memsize`.** Zero dependencies, but spawns a process, parses
  stdout as text, and adds failure modes a syscall cannot have (binary not on PATH).
- **Bare `64 GB` output.** Closest to the literal request, but a later CPU section would force
  a re-decision of the whole output format.
- **`../hw`'s header block.** The only option that makes both repos print identically; declined
  as six lines of ceremony around one number.
- **`1024^3` labeled `GiB`** (what `../hw`'s spec chose) and **`1000^3` labeled `GB`**
  (SI-correct, but would print `68.7 GB` for a Mac everyone calls 64 GB).
- **`stderr` + exit 1 on read failure** (cleanest for a single-purpose tool; matches `../hw`'s
  CPU-block rule) and **`N/A` + exit 1** (stable shape, honest signal). Declined for `../hw`
  parity.
- **Skipping design + ADR.** In this schema `tasks` requires `specs` **and** `adr`, and `adr`
  requires `design`, so the two are effectively inseparable.

## Open questions

- **O1 — deferred by the user: how far the chain goes.** Offered: full
  `proposal -> specs -> design -> adr -> tasks`; stop after `specs`; or record design + adr as
  deliberate conditional skips. Not answered — the reply redirected to "save and stop", so this
  is deferred, not decided. My recommendation stands: **full chain through tasks**, because the
  `gopsutil` commitment earns a repo-level ADR and the schema requires design + adr for tasks.
- **O2 — twin divergence is deliberate.** Shape (`Total Memory: 64 GB` vs `../hw`'s header
  block) and unit (`GB` vs `GiB`) differ from the Linux twin; the error path matches verbatim.
  The two tools will not diff clean at archive time.
- **O3 — test scope, unconfirmed default.** One Go unit test for bytes -> `Total Memory: 64 GB`;
  no acceptance-test harness, since this repo ships python/js runners rather than Go.
- **O4 — API names, assumed.** `func Total() (uint64, error)` in `package ram`, with the
  conversion and the `N/A` text in `main`. Not confirmed.
- **O5 — where the syscall/dependency seam sits for testing.** With `gopsutil` the read is a
  package-level function, so `Total()` cannot be unit-tested without an injected reader. Whether
  to add an injection seam is a `design.md` question, not a clarification.

