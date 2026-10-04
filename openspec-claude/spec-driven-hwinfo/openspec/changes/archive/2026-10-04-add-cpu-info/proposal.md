## Why

`hwinfo` is meant to print basic hardware information, but it currently prints nothing. CPU details (who made it, which model, how fast, how many cores) are the first thing users expect, so they are the natural first capability.

## What Changes

- `hwinfo` prints a human-readable CPU section with vendor, model, speed, and physical/logical core counts.
- Vendor is parsed from the first word of the CPU brand string, with trademark marks such as `(R)` and `(TM)` stripped, because Apple Silicon exposes no vendor value.
- Speed is reported when the OS exposes it. Otherwise it is shown as `unavailable` and the program still exits successfully.
- Core counts are totals only (physical and logical); no per-core-type breakdown.
- macOS only. Apple Silicon is the guaranteed target; Intel Macs are best-effort and untested.
- Out of scope: JSON output, per-core-type breakdown, decoding private IOKit tables for Apple Silicon frequency, live/current speed.

## Capabilities

### New Capabilities
- `cpu-info`: Report CPU vendor, model, speed, and physical/logical core counts as human-readable text on macOS.

### Modified Capabilities

## Impact

- New code in `./internal/cpu` that reads CPU values from macOS sysctl, plus wiring in `main.go`.
- Possible dependency on `golang.org/x/sys/unix` for the 64-bit frequency read (the stdlib `syscall` package does not cover it); to be decided in design.
- No existing specs or APIs are affected.
