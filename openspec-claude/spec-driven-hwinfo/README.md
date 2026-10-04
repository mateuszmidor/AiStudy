# hwinfo

A cli app written in Go that prints out basic information about users hardware.

## Status

- Reports **CPU** info only (vendor, model, speed, core counts).
- Supports **macOS** only. Information is read via `sysctl` (no cgo, no subprocesses). Other platforms are an explicit non-goal for now.
- Values that can't be determined are printed as `unavailable`, and the program always exits with code 0.

## Usage

```
go run .
```

Example output (Apple Silicon):

```
CPU
  Vendor: Apple
  Model: Apple M5 Pro
  Speed: unavailable
  Cores: 18 physical, 18 logical
```

On Apple Silicon the CPU speed is not exposed by the OS, so it is shown as `unavailable`.

## Development

```
go build ./...
go vet ./...
go test ./...
```

## Project layout

- `main.go` – wires the system CPU source into the reader and prints the result.
- `internal/cpu` – CPU reading and formatting logic. The real sysctl-backed source is the only macOS-specific file; everything else is testable on any machine using a fake source.
- `openspec/` – change proposals, designs and specs, managed with [OpenSpec](https://github.com/Fission-AI/OpenSpec) (spec-driven workflow).
