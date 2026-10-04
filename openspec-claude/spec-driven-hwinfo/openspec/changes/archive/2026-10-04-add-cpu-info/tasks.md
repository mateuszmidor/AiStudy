## 1. Setup

- [x] 1.1 Add `golang.org/x/sys` to `go.mod` (`go get golang.org/x/sys`) and verify `go.mod` lists it and `go build ./...` succeeds
- [x] 1.2 Create the `./internal/cpu` package directory (package `cpu`) with a darwin-only sysctl file layout, and verify `go vet ./...` passes on the empty skeleton

## 2. CPU data model and parsing

- [x] 2.1 Define the `source` interface (string, uint32 and uint64 lookups that can fail) and the `CPU` struct (`Vendor`, `Model`, `SpeedHz`, `Physical`, `Logical`, with per-field availability), and verify the package compiles
- [x] 2.2 Implement vendor parsing (first word of brand string, strip `(R)`/`(TM)`/`(tm)`) and verify unit tests pass for `Apple M5 Pro` -> `Apple`, `Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz` -> `Intel`, and empty string -> unavailable
- [x] 2.3 Implement `cpu.Read(source)` that fills the struct from `machdep.cpu.brand_string`, `hw.physicalcpu`, `hw.logicalcpu` and `hw.cpufrequency_max`, marking failed lookups unavailable, and verify fake-source tests cover: 18/18 cores, 6/12 cores, frequency present, frequency missing, empty brand string

## 3. Formatting

- [x] 3.1 Implement the formatter that renders `Vendor`, `Model`, `Speed` and `Cores` lines, showing `unavailable` for missing values, and verify tests assert `Cores: 18 physical, 18 logical` and `Cores: 6 physical, 12 logical`
- [x] 3.2 Implement speed formatting (Hz -> GHz, two decimals) and verify a test asserts `2600000000` -> `2.60 GHz` and a missing frequency -> `unavailable`

## 4. macOS source and CLI wiring

- [x] 4.1 Implement the real `source` in a `//go:build darwin` file using `unix.Sysctl`, `unix.SysctlUint32` and `unix.SysctlUint64`, and verify `GOOS=darwin go build ./...` succeeds
- [x] 4.2 Wire `main.go` (repository root, importing `hwinfo/internal/cpu`) to read the CPU via the real source, print the CPU section, and always exit 0, and verify `go run .` prints `Vendor`, `Model`, `Speed` and `Cores` lines

## 5. Verification

- [x] 5.1 Run `go test ./...` and verify all tests pass
- [x] 5.2 Run `go run .` on this Apple Silicon Mac and verify the output shows `Vendor: Apple`, `Model: Apple M5 Pro`, `Speed: unavailable`, `Cores: 18 physical, 18 logical` and exits with status 0 (`echo $?`)
