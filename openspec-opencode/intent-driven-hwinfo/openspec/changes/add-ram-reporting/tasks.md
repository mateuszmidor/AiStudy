## 1. Implement RAM reporting

- [ ] 1.1 Add `gopsutil/v4/mem` import and call `mem.VirtualMemory()` to read total physical RAM in bytes
- [ ] 1.2 Convert bytes to integer GiB (divide by 1024^3, truncate)
- [ ] 1.3 Add output block with "HWInfo - Memory Information" header and "Total Memory: <value> GiB" line (or "Total Memory: N/A" on error)

## 2. Validate

- [ ] 2.1 Run `go build ./...` and verify compilation
- [ ] 2.2 Run the binary and verify memory line appears with correct GiB value
- [ ] 2.3 Run `openspec validate add-ram-reporting --type change --strict` to validate artifacts
