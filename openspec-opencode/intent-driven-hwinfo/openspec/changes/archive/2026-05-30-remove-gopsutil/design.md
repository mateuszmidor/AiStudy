## Context

The `hwinfo` CLI currently uses two external libraries for CPU data:
- `github.com/shirou/gopsutil/v4/cpu` — for vendor, model, and core counts (reads /proc/cpuinfo)
- `github.com/prometheus/procfs/sysfs` — for frequency data (reads sysfs cpufreq)

The prometheus/procfs root package also provides a `CPUInfo()` function that reads and parses /proc/cpuinfo, covering the same fields gopsutil provides. Since prometheus/procfs is already a direct dependency (used by the sysfs sub-package), gopsutil is redundant.

## Goals / Non-Goals

**Goals:**
- Eliminate the gopsutil dependency entirely
- Preserve identical output for all fields (manufacturer, model, architecture, physical/logical cores)
- Keep the same error behavior: exit 1 if CPU info cannot be read

**Non-Goals:**
- No change to frequency data sourcing (stays on prometheus/procfs/sysfs)
- No change to output format or CLI behavior

## Decisions

1. **Use prometheus/procfs CPUInfo() instead of gopsutil cpu.Info()**
   - Both read `/proc/cpuinfo` and parse the same fields
   - procfs.CPUInfo() provides VendorID, ModelName, PhysicalID, CoreID, CPUCores, Siblings — everything needed
   - procfs uses build tags to select the correct parser: `cpuinfo_x86.go` for linux+amd64, which matches our target

2. **Physical core counting via unique (PhysicalID, CoreID) pairs**
   - gopsutil's `cpu.Counts(false)` internally deduplicates by physical package and core ID
   - We replicate this by iterating parsed CPUInfo entries and building a set of `"physicalID:coreID"` strings
   - On x86, both fields are always present; on other arches where they're empty, fall back to logical count
   - Implementation: ~7 lines, no new dependencies

3. **Logical core counting via len(cpuInfo)**
   - Each entry in the CPUInfo slice corresponds to one logical CPU
   - Direct replacement for `cpu.Counts(true)`

## Risks / Trade-offs

- **[Low] Non-x86 edge case**: On architectures where /proc/cpuinfo lacks physical/core IDs, physical core count falls back to logical count. This matches gopsutil behavior on those platforms.
- **[Low] Different error messages**: procfs returns Go-native errors from /proc/cpuinfo parsing. The existing behavior (stderr message + exit 1) is preserved by checking the error and handling identically.

## Migration Plan

Single atomic change to one file (main.go), verified by running the binary and comparing output. No deployment, migration, or rollback needed — this is a CLI tool.

## Open Questions

- The existing ADR (`adr/0001-use-gopsutil-for-hardware-info.md`) recommends gopsutil. No new ADR is created in this workflow; the existing one is simply no longer in force.
