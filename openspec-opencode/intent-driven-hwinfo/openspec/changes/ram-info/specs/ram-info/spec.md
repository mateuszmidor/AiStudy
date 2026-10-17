# RAM Info

## Purpose

**hwinfo** reports installed **total physical memory** on macOS as a single labeled line on stdout, expressed in **GB**. macOS is the supported platform by intent, not by compiler constraint. The value is sourced via the `gopsutil/v4/mem` package.

## ADDED Requirements

### Requirement: Report total physical memory
The tool SHALL read installed **total physical memory** and print it as a single labeled line on stdout in **GB**.

#### Scenario: Total physical memory displayed as one line
- **GIVEN** the machine reports 68719476736 bytes of **total physical memory**
- **WHEN** the tool runs
- **THEN** it shall print exactly one line to stdout: `Total Memory: 64 GB`
- **AND** it shall print nothing else to stdout
- **AND** it shall exit with status 0

#### Scenario: A total that is not a whole number of GB truncates down
- **GIVEN** the machine reports 17179869183 bytes of **total physical memory**, one byte less than 16 **GB**
- **WHEN** the tool runs
- **THEN** it shall print `Total Memory: 15 GB`
- **AND** the printed value shall be the byte count divided by 1024^3 with the decimal portion truncated, never rounded up

#### Scenario: Total physical memory cannot be read
- **GIVEN** the read of **total physical memory** fails
- **WHEN** the tool runs
- **THEN** it shall print exactly one line to stdout: `Total Memory: N/A`
- **AND** it shall write nothing to stderr
- **AND** it shall exit with status 0

#### Scenario: A read reporting zero bytes is treated as unavailable
- **GIVEN** the read of **total physical memory** succeeds but reports 0 bytes
- **WHEN** the tool runs
- **THEN** it shall print exactly one line to stdout: `Total Memory: N/A`
- **AND** it shall write nothing to stderr
- **AND** it shall exit with status 0

#### Scenario: Arguments and flags are ignored
- **GIVEN** the tool is invoked with any arguments or flags
- **WHEN** the tool runs
- **THEN** it shall ignore all arguments and flags
- **AND** it shall print the total memory line normally

## MODIFIED Requirements

<!-- None. This change introduces a new capability; no existing spec is modified. -->

## REMOVED Requirements

<!-- None. This change removes no behaviour. -->
