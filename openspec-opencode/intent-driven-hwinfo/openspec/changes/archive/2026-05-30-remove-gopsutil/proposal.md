## Why

Remove the dependency on `github.com/shirou/gopsutil/v4/cpu` to reduce the dependency footprint. gopsutil pulls in 8 indirect dependencies (go-ole, wmi, purego, etc.) that are unnecessary on Linux amd64. The prometheus/procfs package is already a direct dependency (used for sysfs cpufreq data) and can fully replace gopsutil for CPU vendor, model, and core count data.

## What Changes

- Remove import of `github.com/shirou/gopsutil/v4/cpu` from `main.go`
- Add import of `github.com/prometheus/procfs` (the root package, not the already-imported sysfs sub-package)
- Replace `cpu.Info()` with `procfs.NewDefaultFS().CPUInfo()` for VendorID and ModelName
- Replace `cpu.Counts(false)` with unique `(PhysicalID, CoreID)` pair dedup for physical core count
- Replace `cpu.Counts(true)` with `len(cpuInfo)` for logical core count
- Run `go mod tidy` to remove gopsutil and its transitive dependencies from go.mod and go.sum

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `cpu-info-display`: Data source for non-frequency CPU information (manufacturer, model, core counts) changes from gopsutil to prometheus/procfs. Behaviour is identical — only the underlying library changes.

## Impact

- `main.go`: import changes and ~7 lines of new dedup logic
- `go.mod` / `go.sum`: removal of `github.com/shirou/gopsutil/v4` and its 8 transitive dependencies
- No change to frequency data sourcing (already using prometheus/procfs/sysfs)
