## Context

Speed values are currently displayed as raw float64 via `fmt.Sprintf("%v", f)`, producing output like "1397.4055624999999". This is overly precise and noisy for a hardware info CLI. Users expect clean integer MHz values.

## Goals / Non-Goals

**Goals:**
- Display CPU Speed (MHz) and Max CPU Speed (MHz) as truncated integer values
- Preserve averaging logic — truncation happens after averaging

**Non-Goals:**
- No rounding — use truncation (floor toward zero) per user preference
- No changes to N/A handling, error paths, or other fields

## Decisions

1. **Truncation via `int()` conversion**: Go's `int(float64)` truncates toward zero, which for positive MHz values is equivalent to `math.Floor`. No need for `math` import.
2. **`fmt.Sprintf("%d", ...)` over `fmt.Printf(...)`**: Replacing `fmt.Sprintf("%v", avgMhz)` with `fmt.Sprintf("%d", int(avgMhz))` for both speed fields. Keeps the existing `avgMhzStr`/`maxMhzStr` string pattern unchanged.
3. **Both fields truncated identically**: Current speed and max speed both get the same treatment for consistency.

## Risks / Trade-offs

- **[Risk] Precision loss**: Users wanting sub-MHz precision lose it. **Accepted** — this is a hardware info tool, not a benchmark. Integer MHz is the standard display format (lscpu, cpufreq-info, etc.).

## Migration Plan

Single line change in `main.go`: replace two `fmt.Sprintf("%v", ...)` calls with `fmt.Sprintf("%d", int(...))`. Rollback is reverting the two format strings.

## Open Questions

None.
