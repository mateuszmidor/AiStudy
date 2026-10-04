## Context

The tool currently reports CPU information (manufacturer, cores, speed, temperature). Total physical RAM is a fundamental hardware metric not yet displayed. The existing `gopsutil/v4` dependency (added for CPU temp) includes a `mem` sub-package that reads total RAM from `/proc/meminfo`.

ADR-0002 removed gopsutil/v3 for bloat, but gopsutil/v4 is already a dependency (added in `add-cpu-temp-reporting`), so using its `mem` sub-package adds zero new dependencies.

## Goals / Non-Goals

**Goals:**
- Display total physical RAM in GiB (integer truncation)
- Best-effort: error does not block existing CPU output

**Non-Goals:**
- Used/available memory
- Memory speed or type
- Module/slot information
- Swap reporting

## Decisions

### Data Source: gopsutil/v4/mem

Use `gopsutil/v4/mem.VirtualMemory()` which reads `/proc/meminfo` and returns `Total` in bytes. Divide by `1024^3` and truncate to integer GiB. Already a transitive dependency — just add the import.

One line:
```go
v, _ := mem.VirtualMemory()
giB := v.Total / (1024 * 1024 * 1024)
```

### Output format

Add a new section after the CPU block to keep output organized:

```
HWInfo - CPU Information
...
CPU Temperature (Avg): 42°C

HWInfo - Memory Information
------------------------
Total Memory: 14 GiB
```

Separate header section follows the existing pattern and groups related RAM info for future expansion (e.g., used memory, DIMM details).

### Error handling

If `mem.VirtualMemory()` fails, display `Total Memory: N/A` and exit 0. Matches the existing best-effort pattern for CPU speed and temperature.

## Risks / Trade-offs

- None significant. The change is minimal, no new deps, well-tested code path.

## Migration Plan

Single-step: add import and memory reading block to main.go.

## Open Questions

- None.
