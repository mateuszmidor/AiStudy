## ADDED Requirements

### Requirement: Display CPU information
The tool SHALL read CPU hardware information and display it in a structured human-readable format.

#### Scenario: Display all CPU fields successfully
- **GIVEN** the system has CPU information accessible via gopsutil
- **WHEN** the tool runs
- **THEN** it shall print a header block "HWInfo - CPU Information\n------------------------"
- **AND** it shall print the following fields in order: CPU Manufacturer, CPU Model, CPU Architecture, Physical Cores, Logical Cores, CPU Speed (MHz), Max CPU Speed (MHz)
- **AND** each field shall be formatted as "<Label>: <Value>" on its own line

#### Scenario: CPU speed values are averaged
- **GIVEN** multiple logical CPUs with different speed values
- **WHEN** the tool displays CPU Speed and Max CPU Speed
- **THEN** the values shall be the arithmetic mean across all logical CPUs/cores
- **AND** the values shall be displayed in MHz without rounding

#### Scenario: Missing single field shows N/A
- **GIVEN** CPU information is partially available (e.g., Max CPU Speed is not provided by the system)
- **WHEN** the tool runs
- **THEN** the unavailable field shall display "N/A" instead of a value
- **AND** the tool shall exit with status 0

#### Scenario: Critical error causes non-zero exit
- **GIVEN** gopsutil fails to read CPU information (e.g., platform not supported)
- **WHEN** the tool runs
- **THEN** it shall print an error message to stderr
- **AND** it shall exit with status 1
- **AND** no partial output shall be printed to stdout

#### Scenario: No arguments or flags accepted
- **GIVEN** the tool is invoked with any arguments or flags
- **WHEN** the tool runs
- **THEN** it shall ignore all arguments and flags
- **AND** display CPU information normally

## MODIFIED Requirements

<!-- No existing specs to modify -->

## REMOVED Requirements

<!-- No existing specs to remove -->
