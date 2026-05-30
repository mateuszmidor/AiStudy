# Memory Info Display

## Purpose

Display total physical RAM in a structured human-readable format. Sourced via the gopsutil/v4/mem Go package reading /proc/meminfo.

## ADDED Requirements

### Requirement: Display total physical memory
The tool SHALL read total physical RAM and display it in GiB.

#### Scenario: Total memory displayed successfully
- **GIVEN** the system has memory information accessible via gopsutil/v4/mem
- **WHEN** the tool runs
- **THEN** it shall print a header block "HWInfo - Memory Information\n------------------------"
- **AND** it shall print "Total Memory: <value> GiB"
- **AND** the value shall be total physical RAM in bytes divided by 1024^3, truncated to integer

#### Scenario: Memory information unavailable shows N/A
- **GIVEN** gopsutil/v4/mem cannot read memory information
- **WHEN** the tool runs
- **THEN** the header block "HWInfo - Memory Information\n------------------------" shall still be printed
- **AND** "Total Memory: N/A" shall be displayed
- **AND** the tool shall exit with status 0

### Requirement: Memory data source
The tool SHALL use gopsutil/v4/mem as the data source for total physical RAM.

#### Scenario: Total RAM from VirtualMemory
- **GIVEN** gopsutil/v4/mem.VirtualMemory() succeeds
- **WHEN** the tool reads total memory
- **THEN** the value shall be derived from the Total field (in bytes), divided by 1024^3
