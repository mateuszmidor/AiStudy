## Purpose

Lets users see basic facts about the CPU of the machine they run `hwinfo` on: who makes it, which model it is, how fast it runs, and how many cores it has.

## Requirements

### Requirement: CPU section is printed as labeled text
The system SHALL print a CPU section as human-readable text, with one labeled line each for vendor, model, speed, and cores.

#### Scenario: Running on a supported Mac
- **WHEN** the user runs `hwinfo` on macOS
- **THEN** the output contains a CPU section with `Vendor`, `Model`, `Speed`, and `Cores` lines
- **AND** the program exits with status 0

### Requirement: CPU model is reported
The system SHALL report the CPU model exactly as the operating system exposes its brand string.

#### Scenario: Apple Silicon model
- **WHEN** the OS reports the brand string `Apple M5 Pro`
- **THEN** the `Model` line shows `Apple M5 Pro`

### Requirement: CPU vendor is derived from the model
The system SHALL report the CPU vendor as the first word of the model (brand string), with trademark marks such as `(R)` and `(TM)` removed, because not every CPU exposes a vendor value.

#### Scenario: Apple Silicon vendor
- **WHEN** the model is `Apple M5 Pro`
- **THEN** the `Vendor` line shows `Apple`

#### Scenario: Intel vendor with trademark marks
- **WHEN** the model is `Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz`
- **THEN** the `Vendor` line shows `Intel`

#### Scenario: Brand string is empty
- **WHEN** the OS reports an empty brand string
- **THEN** the `Model` and `Vendor` lines show `unavailable`
- **AND** the program exits with status 0

### Requirement: Core counts are reported as totals
The system SHALL report the number of physical cores and the number of logical cores as totals, without a per-core-type breakdown.

#### Scenario: Chip without simultaneous multithreading
- **WHEN** the machine has 18 physical cores and 18 logical cores
- **THEN** the `Cores` line shows `18 physical, 18 logical`

#### Scenario: Chip with more logical than physical cores
- **WHEN** the machine has 6 physical cores and 12 logical cores
- **THEN** the `Cores` line shows `6 physical, 12 logical`

### Requirement: CPU speed is reported when available
The system SHALL report the maximum CPU speed when the operating system exposes it, and SHALL show `unavailable` otherwise, without treating the absence as an error.

#### Scenario: Speed exposed by the OS
- **WHEN** the OS exposes a maximum CPU frequency of 2600000000 Hz
- **THEN** the `Speed` line shows `2.60 GHz`

#### Scenario: Speed not exposed by the OS
- **WHEN** the OS exposes no CPU frequency (as on Apple Silicon)
- **THEN** the `Speed` line shows `unavailable`
- **AND** the program exits with status 0

### Requirement: Platform support
The system SHALL support Apple Silicon Macs. Intel Macs SHOULD work on a best-effort basis but are not guaranteed.

#### Scenario: Apple Silicon Mac
- **WHEN** the user runs `hwinfo` on an Apple Silicon Mac
- **THEN** every line of the CPU section is either a value or `unavailable`
