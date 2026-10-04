## Context

`hwinfo` started as an empty Go module (`main.go` contained only `package main`, no dependencies). See proposal.md for motivation and scope.

Findings from probing an Apple Silicon Mac (arm64):
- `machdep.cpu.brand_string`, `hw.physicalcpu` and `hw.logicalcpu` exist.
- `machdep.cpu.vendor` and `hw.cpufrequency(_max)` do not exist (sysctl reports "unknown oid").
- A missing sysctl is therefore a normal, expected condition, not a failure.

## Goals / Non-Goals

**Goals:**
- Read CPU facts from sysctl without cgo or subprocesses.
- Keep the logic that turns raw values into display text (vendor parsing, GHz formatting, `unavailable` fallback) testable without a real Mac.

**Non-Goals:**
- Any non-macOS code path. Other platforms are not built or supported.
- Decoding private IOKit/ioreg data to get an Apple Silicon frequency.

## Decisions

### Read sysctl through `golang.org/x/sys/unix`
Use `unix.Sysctl` (strings), `unix.SysctlUint32` (core counts) and `unix.SysctlUint64` (frequency).

- Alternative: shell out to `sysctl -n`. Rejected: a subprocess per value and text parsing for no benefit.
- Alternative: stdlib `syscall.Sysctl`. Rejected: it returns a string and trims a trailing NUL byte, so a raw 64-bit value cannot be read reliably, and the stdlib has no uint64 variant.
- Cost: one new module dependency (`golang.org/x/sys`), which the Go team maintains.

### Separate reading from formatting
Split into three small pieces:
1. An exported `Source` interface (string, uint32 and uint64 lookups that return an error when a value does not exist) with a real implementation, `SystemSource`, backed by `x/sys/unix`. It is exported so `main.go` can pass it to `cpu.Read`.
2. A `cpu` package (`internal/cpu`) whose `Read(Source)` asks a source for values and returns a `CPU` struct (`Vendor`, `Model`, `SpeedHz`, `Physical`, `Logical`). A zero value means the fact is unavailable (a CPU with 0 cores or 0 Hz is not meaningful), so the struct carries no separate availability flags.
3. A `Format(CPU)` function that renders the struct as labeled lines.

Tests feed a fake source, so every scenario in the spec (Intel brand string, empty brand string, 6/12 cores, 2.6 GHz, missing frequency) runs on any machine, including this Apple Silicon Mac. This is also how Intel behavior is covered without an Intel Mac.

- Alternative: one function in `main.go` calling sysctl directly. Rejected: simpler, but the Intel and empty-string scenarios would be untestable.

### Package location
All CPU code lives in `./internal/cpu` (import path `hwinfo/internal/cpu`, per the module name in `go.mod`), so it is private to this module. `main.go` stays at the repository root and only wires the package to the real source and prints the result.

### Build constraint for macOS
Put the real sysctl implementation in a file with a `//go:build darwin` constraint. The cpu logic and formatter stay platform-neutral so their tests run anywhere.

### Missing values are `unavailable`, never an error
A sysctl lookup that fails (unknown oid) leaves the field at its zero value, meaning unavailable. Only the formatter turns that into the text `unavailable`, including the `Cores` line when either core count is missing. `main` exits 0 regardless.

### Vendor parsing
Trim surrounding whitespace from the brand string to get the model, then take its first whitespace-separated word and remove `(R)`, `(TM)` and `(tm)` to get the vendor. An empty brand string produces an unavailable vendor and model.

### Speed formatting
Convert Hz to GHz with two decimals (`2600000000` becomes `2.60 GHz`).

## Risks / Trade-offs

- [Frequency code path cannot be exercised on real hardware here (Apple Silicon has no such sysctl)] → Covered by fake-source tests; documented as best-effort in the spec.
- [First-word vendor parsing misreads unusual brand strings] → Limited to the unavailable/empty case and the known Apple and Intel formats; revisit if other brand strings appear.
- [New dependency `golang.org/x/sys`] → Small, widely used and maintained by the Go project; pin through `go.mod`.
- [Apple may rename or drop sysctl keys in future macOS versions] → Missing keys already degrade to `unavailable`.
