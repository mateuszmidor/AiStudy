## Context

The project currently uses `prometheus/procfs` and `prometheus/procfs/sysfs` for CPU info and frequency data. ADR-0002 removed `gopsutil/v3` due to transitive dep bloat, but this change adds `gopsutil/v4/sensors` as a focused dependency for hardware temperature reading.

CPU temperature is available via the `gopsutil/v4/sensors` package, which reads from Linux hwmon/thermal_zone (and has cross-platform support). Sensor keys follow the pattern `k10temp_tctl`, `coretemp_package`, etc.

## Goals / Non-Goals

**Goals:**
- Read CPU die/package temperature and display average across CPU sensors
- Best-effort: temperature failure does not block other output
- Use a Go package rather than direct sysfs file reads

**Non-Goals:**
- Per-core temperature display (user requested average)
- GPU, NVMe, or other non-CPU sensor display
- Continuous monitoring or delta tracking

## Decisions

### Data Source: gopsutil/v4/sensors

**Chosen**: `github.com/shirou/gopsutil/v4/sensors` — the de facto standard Go library for hardware sensor data. Returns labeled temperature stats (e.g., `k10temp_tctl`, `coretemp_package_0`) as float64 Celsius values.

**Alternatives considered**:

| Option | Deps | Rationale for rejection |
|--------|------|------------------------|
| Direct hwmon sysfs parsing | 0 | User requested a Go package rather than direct filesystem access |
| `northwindlight/cputemp` | ~4 | Only reads thermal_zone, misses hwmon/k10temp on AMD |
| `prometheus/procfs/sysfs.ClassThermalZoneStats()` | 0 | Only reads ACPI thermal zones, not CPU-die hwmon sensors |

### CPU sensor filtering

Filter gopsutil results by matching sensor keys containing `k10temp` (AMD), `coretemp` (Intel), or `zen` (newer AMD). This excludes non-CPU sensors like `acpitz`, `amdgpu`, `nvme`.

### Averaging and error handling

Average matched sensor values (e.g., Tdie + TccdX on multi-CCD AMD) for a single average. If no CPU sensors found or any error occurs, display N/A. Matches existing pattern for CPU speed.

## Risks / Trade-offs

- **[ADR-0002 tension]** ADR-0002 removed gopsutil for excess deps. This change adds a subset (`sensors`) which adds ~4 transitive deps (purego, go-ole, wmi, yusufpapurcu/wmi). Acceptable tradeoff for using a maintained, tested library.
- **[Virtual machines / containers]** Sensor data is often absent. Handled by N/A fallback.

## Migration Plan

Single-step: add gopsutil/v4/sensors dependency and replace direct hwmon reading with the library call.

## Open Questions

- None.
