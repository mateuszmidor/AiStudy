## Context

The current implementation uses `gopsutil/v4/cpu`'s `Info()` call for all CPU data. On Linux, gopsutil overrides the `Mhz` field with the value from `/sys/devices/system/cpu/cpuN/cpufreq/cpuinfo_max_freq` — it only exposes the **max** frequency, not the current/actual speed. The tool's "CPU Speed (MHz)" therefore shows max speed mislabeled, and "Max CPU Speed (MHz)" is hardcoded to `"N/A"`.

The kernel exposes both current and max frequency via distinct sysfs files (`scaling_cur_freq`, `cpuinfo_max_freq`). gopsutil reads only `cpuinfo_max_freq` and discards the current speed. No other gopsutil API provides current frequency.

## Goals / Non-Goals

**Goals:**
- Populate "CPU Speed (MHz)" with actual current frequency (mean across logical CPUs)
- Populate "Max CPU Speed (MHz)" with hardware max frequency (mean across logical CPUs)
- Keep gopsutil for manufacturer, model, architecture, physical/logical core counts
- Maintain all existing output format, error handling, and N/A fallback behavior
- Handle missing cpufreq gracefully (show N/A, exit 0)

**Non-Goals:**
- No output format changes (same labels, same ordering)
- No changes to non-frequency fields
- No cross-platform changes (Linux amd64 only, matching existing scope)

## Decisions

1. **`prometheus/procfs/sysfs` over direct sysfs reads**: The package provides `SystemCpufreq()` which returns a structured `[]SystemCPUCpufreqStats` with `ScalingCurrentFrequency` and `CpuinfoMaximumFrequency` per CPU. Compared to manual `os.ReadFile` + glob, it handles:
   - Error classification (permission vs not-exist vs unknown)
   - Parallel reads via errgroup
   - Consistent parsing from kHz strings to uint64

2. **`ScalingCurrentFrequency` over `cpuinfo_cur_freq`**: `scaling_cur_freq` reports the frequency chosen by the cpufreq governor, which matches what the OS is actively requesting. `cpuinfo_cur_freq` reports the raw hardware frequency, which may be stale on some drivers. `scaling_cur_freq` is more widely available and semantically correct for "current CPU speed".

3. **Averaging across logical CPUs**: Same approach as current code — compute arithmetic mean across all logical CPUs. This matches the existing spec and preserves consistency with the current averaging behavior.

4. **kHz→MHz conversion**: sysfs values are in kHz. Divide by 1000 for MHz display. Use `float64` arithmetic to preserve precision, matching the existing "no rounding" behavior.

5. **gopsutil retained for non-frequency fields**: No reason to change manufacturer, model, architecture, or core count sourcing. gopsutil handles these well.

## Risks / Trade-offs

- **[Risk] Missing cpufreq sysfs**: Some environments (containers, VMs, very old kernels) lack cpufreq support → `SystemCpufreq()` returns error → both speed fields show N/A. **Mitigation**: Error is caught, both fields display "N/A", tool exits 0 (existing pattern).
- **[Risk] Partial cpufreq (some CPUs, not others)**: Unusual but possible if CPUs are hotplugged or heterogeneous. **Mitigation**: The prometheus package handles this gracefully — missing CPUs are skipped in the slice. Only available CPUs contribute to the average.
- **[Trade-off] New dependency**: Adds `github.com/prometheus/procfs` (+15 imports from stdlib). Already in Go module cache at v0.10.1. Well-maintained (870+ stars, Prometheus project).
- **[Risk] Prometheus procfs API stability**: The package is v0.x (pre-1.0). API changes possible but unlikely for the cpufreq subsystem. Mitigated by go.mod pinning.

## Migration Plan

1. Add `github.com/prometheus/procfs` to go.mod
2. Rewrite frequency computation in main.go:
   - Remove `totalMhz`/`count` loop over `cpu.Info()` for frequency
   - Add `sysfs.SystemCpufreq()` call alongside existing `cpu.Info()`
   - Compute `currentSpeed` from `ScalingCurrentFrequency` values (kHz→MHz)
   - Compute `maxSpeed` from `CpuinfoMaximumFrequency` values (kHz→MHz)
   - Fall back to "N/A" for either field on error
3. Update spec to reflect dual data source
4. Build and verify against real hardware

## Open Questions

None — all design decisions are resolved.
