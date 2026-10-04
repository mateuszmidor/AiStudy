## 1. Project scaffold

- [x] 1.1 Create `go.mod` at repo root targeting Go 1.26, adding `github.com/shirou/gopsutil/v4/cpu` dependency
- [x] 1.2 Create `main.go` at repo root with package `main` and imports of gopsutil/cpu, fmt, and os

## 2. Core CPU data reading

- [x] 2.1 Implement CPU info retrieval via `cpu.Info()` — extract manufacturer (VendorID), model (ModelName), architecture, physical/logical core counts, and per-core speed values
- [x] 2.2 Implement CPU speed averaging: compute mean of per-core MHz values across all logical CPUs for both current speed and max speed
- [x] 2.3 Handle nil/empty cpu.Info() result as critical error (exit 1 with stderr message)
- [x] 2.4 Handle individual missing fields gracefully (display `N/A` for that field)

## 3. Output formatting

- [x] 3.1 Print header block: "HWInfo - CPU Information" followed by "------------------------"
- [x] 3.2 Print all 7 fields in strict order with "<Label>: <Value>" format: CPU Manufacturer, CPU Model, CPU Architecture, Physical Cores, Logical Cores, CPU Speed (MHz), Max CPU Speed (MHz)
- [x] 3.3 Run `go build` and verify the tool compiles successfully
- [x] 3.4 Run `openspec validate hwinfo-cpu-cli --type change --strict` to verify artifacts pass validation
