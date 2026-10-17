# Use gopsutil/v4/mem for reading total physical memory

## Status

Accepted

Date: 2026-10-13

## Context and Problem Statement

**hwinfo** needs to report installed **total physical memory** on macOS, and this change is the tool's first slice. macOS itself exposes the value through the `hw.memsize` sysctl, so a direct syscall would be enough for the number itself. The deciding context is what comes next: the repository intends further memory reporting, the data source is the one choice that outlives a single change, and the Linux twin `../hw` already reads memory with a library rather than a direct syscall. This repository had no dependencies before this change.

## Considered Options

- **`gopsutil/v4/mem`** — `mem.VirtualMemory().Total` returns the total in bytes, and the same call also yields `Available`, `Used`, `Free`, `Wired`, and `UsedPercent`. On darwin the total is ultimately the `hw.memsize` sysctl, reached through the library's mach VM plumbing (`common.NewSystemLib`, `HostStatistics`, and a `Dlsym("vm_kernel_page_size")`), which brings `ebitengine/purego` and other transitive modules with it.
- **`golang.org/x/sys/unix`** — `unix.SysctlUint64("hw.memsize")` gets the same number with the fewest moving parts: one direct dependency, no purego, no runtime `dlopen`, roughly five lines, and the module was already in the local cache. It returns only the total, so any later memory-details work would need a second data source, and it diverges from how the twin reads memory.
- **`os/exec` with `sysctl -n hw.memsize`** — no dependency at all, but it spawns a process, parses stdout as text, and cannot fail the way a syscall cannot: the failure set grows to include "binary not on `PATH`".

## Decision Outcome

Chosen option: "`gopsutil/v4/mem`", because it yields the total today and the remaining memory statistics from the same call tomorrow, so a future memory-details slice is a wider print rather than a new data source; and because it matches the twin's approach, so both tools read memory the same way. The version is pinned in `go.mod` at `v4.26.9`, which was already present in the local module cache.

### Consequences

- Good, because one dependency covers the total now and the other memory statistics later, which is why the heavier option won over the five-line syscall.
- Good, because on darwin `.Total` is the `hw.memsize` sysctl, so the mach VM plumbing runs alongside the read without computing the requested number.
- Good, because the twin reads memory identically, keeping the two implementations comparable even though their output shapes deliberately differ.
- Bad, because it pulls `ebitengine/purego` — a runtime `dlopen` of libSystem — and further transitive modules, in exchange for a single number in this slice.
- Bad, because `.Total` is only reachable through the call that also computes `Available`, `Used`, `Free`, and `Wired`, so the read does more work than the one value requires.
- Bad, because `gopsutil/v4` is comparatively new and its API may change; the exposure is bounded by the `go.mod` pin.
- Neutral, because the output contract — a single labeled line, `GB` meaning 1024³, and in-band `N/A` at exit 0 — is deliberately not recorded here; it is a product convention carried by the merged spec and the glossary.
- Neutral, because the library's runtime `dlopen` means a failure surfaces as an error from the call rather than as an absent value, which is why the caller only ever has to handle a returned error.
