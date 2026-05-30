# Use prometheus/procfs for CPU hardware information

## Context and Problem Statement

The hwinfo CLI needs to read CPU hardware information (manufacturer, model, architecture, core counts) on Linux amd64. The existing dependency `github.com/shirou/gopsutil/v4/cpu` provides this via `/proc/cpuinfo` but pulls in 8 transitive dependencies (go-ole, wmi, purego, etc.) that are unnecessary on Linux. The `github.com/prometheus/procfs` package is already a direct dependency (used for sysfs cpufreq data) and provides equivalent `CPUInfo()` access to `/proc/cpuinfo`.

## Considered Options

- **prometheus/procfs root package** — Already a direct dependency. Provides `CPUInfo()` with VendorID, ModelName, PhysicalID, CoreID, CPUCores — everything gopsutil provides. Saves 8 transitive deps.
- **Raw /proc/cpuinfo parsing** — No external dependency, but format varies across kernel versions. Rejected because procfs already handles this.
- **Keep gopsutil** — Works, but introduces unnecessary transitive dependencies and duplicates functionality already provided by an existing dependency.

## Decision Outcome

Chosen option: "prometheus/procfs root package", because it is already a dependency, covers all needed fields, and removing gopsutil eliminates 8 transitive dependencies with no loss of functionality on Linux amd64.

### Consequences

- Good, because dependency count drops — gopsutil and its 8 transitive deps are removed
- Good, because the library is purpose-built for reading Linux /proc filesystem data with structured types
- Bad, because physical core counting requires a short manual dedup loop instead of a library call — mitigated by the loop being 7 lines of straightforward map logic
