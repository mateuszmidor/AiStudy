## Why

CPU speed values are displayed as full-precision floating point (e.g., "1397.4055624999999"), which is noisy and impractical for a hardware info tool. Users expect clean integer MHz values. Both current speed and max speed should be truncated to integer MHz for readability.

## What Changes

- Change CPU Speed (MHz) display from raw float to truncated integer
- Change Max CPU Speed (MHz) display from raw float to truncated integer
- Update spec to reflect truncation instead of raw float display

## Capabilities

### New Capabilities
*(none)*

### Modified Capabilities
- `cpu-info-display`: Scenario "CPU speed values are averaged" — change "without rounding" to "truncated to integer MHz". All speed values displayed as whole integers.

## Impact

- **Code**: `main.go` only — both `fmt.Sprintf("%v", ...)` calls for speed values change to `fmt.Sprintf("%d", int(...))` for truncation
- **Spec**: One line change in `specs/cpu-info-display/spec.md`
- **No new dependencies**
