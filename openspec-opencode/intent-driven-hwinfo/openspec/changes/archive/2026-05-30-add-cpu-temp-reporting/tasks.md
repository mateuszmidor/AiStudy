## 1. Implement temperature reader using gopsutil/v4/sensors

- [x] 1.1 Add gopsutil/v4/sensors dependency via `go get`
- [x] 1.2 Create function using gopsutil/v4/sensors to fetch all temperature sensors
- [x] 1.3 Filter for CPU sensors by matching keys containing "k10temp", "coretemp", or "zen"
- [x] 1.4 Compute average across matched CPU sensors and convert to integer Celsius

## 2. Integrate into main.go

- [x] 2.1 Call the temperature function in the main flow, after the existing speed averaging block
- [x] 2.2 Add `CPU Temperature (Avg): <value>°C` output line after the existing speed lines, or `CPU Temperature (Avg): N/A` on failure

## 3. Validate

- [x] 3.1 Run `go build ./...` and verify compilation
- [x] 3.2 Run the binary and verify temperature line appears with a sensible value
- [x] 3.3 Run `openspec validate add-cpu-temp-reporting --type change --strict` to validate artifacts
