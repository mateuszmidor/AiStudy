# Use gopsutil/v4/cpu for CPU hardware information

## Context and Problem Statement

The hwinfo CLI needs to read CPU hardware information (manufacturer, model, architecture, core counts, speeds) on Linux amd64. The data is available via `/proc/cpuinfo` but requires manual parsing of unstructured text with varying formats across kernel versions.

## Considered Options

- **gopsutil/v4/cpu**: Stable Go library providing a structured API abstracting over `/proc/cpuinfo`. Used widely in the Go ecosystem.
- **Raw procfs parsing**: Read `/proc/cpuinfo` directly. No external dependency but brittle — format varies across kernel versions and distributions.
- **Manual syscall approach**: Use `unix.Sysctl` or similar low-level calls. More complex, platform-specific, and no benefit over gopsutil.

## Decision Outcome

Chosen option: "gopsutil/v4/cpu", because it provides a tested, stable abstraction over procfs with a clean Go API, avoiding brittle string parsing while adding only one dependency.

### Consequences

- Good, because gopsutil handles cross-kernel-version compatibility internally and abstracts per-CPU data iteration
- Good, because the library is widely used and maintained, reducing bus-factor risk
- Bad, because it introduces an external dependency that must be kept up-to-date
- Bad, because gopsutil v4 is relatively new; API changes are possible but mitigated by go.mod pinning
