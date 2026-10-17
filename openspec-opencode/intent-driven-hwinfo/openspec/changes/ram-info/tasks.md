## 1. Dependency setup

- [x] 1.1 Add `github.com/shirou/gopsutil/v4` pinned at `v4.26.9` with `go get github.com/shirou/gopsutil/v4@v4.26.9`, honoring ADR-0001 — done when `go.mod` requires exactly that version
- [x] 1.2 Run `go mod tidy` to record the transitive requirements and create `go.sum`, and confirm the install did not need the network (use `GOSUMDB=off` for this step only if the checksum database is unreachable)

## 2. Read and format in `internal/ram`

- [x] 2.1 Create `internal/ram/ram.go` with `func Total() (uint64, error)`, reading the total through an unexported `readTotal` function variable that delegates to `mem.VirtualMemory().Total`
- [x] 2.2 Add `func FormatTotalLine(total uint64, available bool) string` to the same package, returning `Total Memory: <n> GB` for bytes divided by 1024³ and truncated, or `Total Memory: N/A` when unavailable
- [x] 2.3 Create `internal/ram/ram_test.go` covering 68719476736 → `Total Memory: 64 GB`, 17179869183 → `Total Memory: 15 GB` (truncates, never rounds up), the unavailable form, and error propagation through the injected `readTotal`
- [x] 2.4 Run `go test ./...` and confirm every test passes

## 3. CLI entry point

- [x] 3.1 Create `main.go` (package `main`) that calls `ram.Total()`, computes availability as `err == nil && total > 0`, prints exactly one line via `ram.FormatTotalLine`, and ignores `os.Args`
- [x] 3.2 Confirm no code path writes to stderr and none exits non-zero — no `os.Exit` call, no logging, no error text
- [x] 3.3 Run `go build ./...` and confirm the binary builds

## 4. Verification on this machine

- [x] 4.1 Run the binary and confirm stdout is exactly one line, `Total Memory: 64 GB`, with nothing else printed
- [x] 4.2 Cross-check the printed value against `sysctl -n hw.memsize` (68719476736 bytes) and About This Mac (64 GB)
- [x] 4.3 Confirm the unavailable form is covered by the unit test, and record in the apply log that the branch cannot be triggered on a healthy machine

## 5. Artifact validation and ADR compliance

- [x] 5.1 Run `openspec validate ram-info --type change --strict` and confirm the change is valid
- [x] 5.2 Confirm the implementation honors ADR-0001 — the read goes through `gopsutil/v4/mem`, with no direct `sysctl` syscall and no `os/exec` fallback
- [x] 5.3 Confirm the change is ready to archive, remembering the config rule that OpenSpec Sync must run before archiving
