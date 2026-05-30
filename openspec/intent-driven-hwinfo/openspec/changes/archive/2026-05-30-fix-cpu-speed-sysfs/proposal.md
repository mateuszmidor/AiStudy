## Why

The tool has two bugs in CPU frequency display: (1) "CPU Speed (MHz)" shows the **max** frequency instead of current/actual speed because gopsutil overrides its `Mhz` field with `cpuinfo_max_freq` on Linux, and (2) "Max CPU Speed (MHz)" is hardcoded to `"N/A"` and never populated. Users get mislabeled data and see "N/A" for a value readily available from sysfs.

## What Changes

- Add `github.com/prometheus/procfs/sysfs` dependency for reading per-CPU frequency from sysfs
- Replace gopsutil-based frequency averaging with `prometheus/procfs/sysfs`'s `SystemCpufreq()` API
- **CPU Speed (MHz)**: computed from `ScalingCurrentFrequency` (actual current speed per core, in kHz)
- **Max CPU Speed (MHz)**: computed from `CpuinfoMaximumFrequency` (hardware max per core, in kHz)
- Keep gopsutil for non-frequency fields (manufacturer, model, core counts)
- Update spec to reflect the new data source for frequency fields

## Capabilities

### New Capabilities
*(none)*

### Modified Capabilities
- `cpu-info-display`: Frequency data source changes from gopsutil to `prometheus/procfs/sysfs` for CPU speed and max CPU speed fields. The averaging and N/A fallback behaviours remain unchanged.

## Impact

- **New dependency**: `github.com/prometheus/procfs` (v0.10.1 already in local module cache)
- **Code**: `main.go` lines 25-42 (frequency computation logic) rewritten
- **Spec**: `specs/cpu-info-display/spec.md` updates to Purpose section (data source) and scenario wording
- **No breaking changes**: Output format and field labels remain identical
