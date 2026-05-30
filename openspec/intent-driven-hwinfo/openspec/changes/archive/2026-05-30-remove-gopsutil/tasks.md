## 1. Replace gopsutil dependency with prometheus/procfs

- [x] 1.1 Update imports in main.go: remove `github.com/shirou/gopsutil/v4/cpu`, add `github.com/prometheus/procfs`
- [x] 1.2 Replace `cpu.Info()` call with `procfs.NewDefaultFS()` + `.CPUInfo()` for VendorID and ModelName
- [x] 1.3 Replace `cpu.Counts(true)` with `len(cpuInfo)` for logical core count
- [x] 1.4 Add unique `(PhysicalID, CoreID)` pair dedup loop to replace `cpu.Counts(false)` for physical core count
- [x] 1.5 Run `go mod tidy` to remove gopsutil and its transitive dependencies from go.mod / go.sum

## 2. Verify

- [x] 2.1 `go build` and `go vet` pass with no errors
- [x] 2.2 Run binary and confirm output matches expected format and values
- [x] 2.3 Run `openspec validate remove-gopsutil --type change --strict`
