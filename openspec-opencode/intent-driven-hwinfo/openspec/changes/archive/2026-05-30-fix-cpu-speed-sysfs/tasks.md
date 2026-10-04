## 1. Add dependency and refactor frequency logic

- [x] 1.1 Add `github.com/prometheus/procfs` to go.mod and run `go mod tidy`
- [x] 1.2 Import `github.com/prometheus/procfs/sysfs` in main.go
- [x] 1.3 Replace gopsutil-based frequency loop with `sysfs.SystemCpufreq()` call
- [x] 1.4 Compute current speed from `ScalingCurrentFrequency` values (kHz → MHz, arithmetic mean)
- [x] 1.5 Compute max speed from `CpuinfoMaximumFrequency` values (kHz → MHz, arithmetic mean)
- [x] 1.6 Handle error from `SystemCpufreq()` — set both speeds to "N/A" without exiting

## 2. Verify and validate

- [x] 2.1 Run `go build ./...` to verify compilation
- [x] 2.2 Run `go run main.go` and confirm both CPU Speed and Max CPU Speed display correct values
- [x] 2.3 Confirm missing cpufreq (e.g., inside a container without /sys access) shows N/A gracefully
