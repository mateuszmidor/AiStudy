## Why

The tool currently reports CPU manufacturer, model, architecture, core counts, and speed — but not temperature. CPU temperature is a fundamental health metric for diagnostics, thermal monitoring, and system assessment. Adding it fills a natural gap in the hardware info display.

## What Changes

- Add CPU temperature reporting to the tool output, displaying average CPU die temperature in Celsius
- Add `github.com/shirou/gopsutil/v4/sensors` as a dependency for cross-platform temperature reading via hwmon/thermal_zone
- Add a new output line: `CPU Temperature (Avg): <value>°C` (or `N/A` if unavailable)
- Temperature is averaged across all CPU temperature sensors (e.g., Tdie + Tccd1 + Tccd2 on AMD, or per-core temps on Intel) into a single average value
- Temperature sourcing is best-effort (not a hard failure) — all other fields still display if temp is unavailable

## Capabilities

### New Capabilities
- `cpu-temp`: Read CPU die/package temperature from hardware sensors via gopsutil, average multiple CPU sensor readings, convert millidegrees to Celsius for display

### Modified Capabilities
- `cpu-info-display`: Add "CPU Temperature (Avg)" as an additional output field after the existing speed fields

## Impact

- **New dependency**: `github.com/shirou/gopsutil/v4/sensors` (~10-15 transitive deps)
- **Modified file**: `main.go` — add temperature reading logic and output line
- **No breaking changes**: All existing output fields and behavior preserved
- **Cross-platform**: gopsutil supports Linux, macOS, and Windows; temp will work on each where sensors are available
