## Why

The tool reports CPU information (manufacturer, cores, speed, temperature) but omits one of the most fundamental hardware metrics: total physical memory. Knowing total RAM is essential for system profiling and diagnostics.

## What Changes

- Add total physical RAM reporting, displayed in GiB (integer truncation)
- Use the existing `gopsutil/v4/mem` dependency (already pulled in by the CPU temp change) to read memory info from `/proc/meminfo`
- Add a new output section "HWInfo - Memory Information" with a "Total Memory" line
- Best-effort: if memory info is unavailable, display `N/A` and exit 0

## Capabilities

### New Capabilities
- `memory-info-display`: Read total physical RAM from the system and display it in GiB with integer truncation

### Modified Capabilities
- None

## Impact

- **No new dependencies**: `gopsutil/v4/mem` is a sub-package of the already-added `gopsutil/v4`
- **Modified file**: `main.go` — add memory reading and output block
- **No breaking changes**: All existing CPU output fields preserved
