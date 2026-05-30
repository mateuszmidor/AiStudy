## Why

Linux users frequently need quick access to CPU hardware information for diagnostics, system profiling, and hardware assessment. Existing tools like `lscpu` require parsing unstructured output, while procfs-based approaches need manual field mapping. A lightweight Go CLI that surfaces essential CPU metrics via a clean library API would provide a fast, portable, and dependency-free inspection tool.

## What Changes

- Create a single-file Go CLI application at repository root (`main.go`)
- Implement CPU information display: manufacturer, model, architecture, core counts, and speed
- Add `github.com/shirou/gopsutil/v4/cpu` as the sole external dependency
- Output human-readable plain text with strict field ordering and a header block
- Handle missing data gracefully (display `N/A`) and critical errors with non-zero exit

## Capabilities

### New Capabilities
- `cpu-info-display`: Read CPU hardware information via gopsutil and display it in a structured human-readable format with fixed labels, strict field ordering, and a header block.

### Modified Capabilities

None — this is a new project with no existing capabilities.

## Impact

- **New code**: `main.go` at repo root — single-file CLI (approx. 80–120 lines)
- **New dependency**: `github.com/shirou/gopsutil/v4/cpu` in `go.mod`
- **No changes** to existing code, APIs, or systems
