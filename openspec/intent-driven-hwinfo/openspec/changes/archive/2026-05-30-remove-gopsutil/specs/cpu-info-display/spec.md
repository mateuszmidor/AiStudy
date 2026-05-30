## MODIFIED Requirements

### Requirement: Display CPU information
The tool SHALL read CPU hardware information and display it in a structured human-readable format.

#### Scenario: Display all CPU fields successfully
- **GIVEN** the system has CPU information accessible via prometheus/procfs
- **WHEN** the tool runs
- **THEN** it shall print a header block "HWInfo - CPU Information\n------------------------"
- **AND** it shall print the following fields in order: CPU Manufacturer, CPU Model, CPU Architecture, Physical Cores, Logical Cores, CPU Speed (MHz), Max CPU Speed (MHz)
- **AND** each field shall be formatted as "<Label>: <Value>" on its own line

#### Scenario: CPU speed values are averaged
- **GIVEN** multiple logical CPUs with different speed values
- **WHEN** the tool displays CPU Speed and Max CPU Speed
- **THEN** the values shall be the arithmetic mean across all logical CPUs/cores
- **AND** the values shall be displayed as integer MHz (decimal portion truncated)

#### Scenario: Missing single field shows N/A
- **GIVEN** CPU information is partially available (e.g., cpufreq sysfs not available, so frequency data cannot be read)
- **WHEN** the tool runs
- **THEN** the unavailable field shall display "N/A" instead of a value
- **AND** the tool shall exit with status 0

#### Scenario: Critical error causes non-zero exit
- **GIVEN** prometheus/procfs fails to read CPU information (e.g., /proc/cpuinfo not accessible)
- **WHEN** the tool runs
- **THEN** it shall print an error message to stderr
- **AND** it shall exit with status 1
- **AND** no partial output shall be printed to stdout

#### Scenario: No arguments or flags accepted
- **GIVEN** the tool is invoked with any arguments or flags
- **WHEN** the tool runs
- **THEN** it shall ignore all arguments and flags
- **AND** display CPU information normally

### Requirement: CPU Speed data source
The tool SHALL use prometheus/procfs/sysfs as the data source for CPU frequency values (current speed and max speed).

#### Scenario: Current CPU Speed from scaling_cur_freq
- **GIVEN** the system has cpufreq enabled with scaling_cur_freq available
- **WHEN** the tool reads CPU Speed
- **THEN** the value shall be the arithmetic mean of scaling_cur_freq values (converted from kHz to MHz) across all logical CPUs

#### Scenario: Max CPU Speed from cpuinfo_max_freq
- **GIVEN** the system has cpufreq enabled with cpuinfo_max_freq available
- **WHEN** the tool reads Max CPU Speed
- **THEN** the value shall be the arithmetic mean of cpuinfo_max_freq values (converted from kHz to MHz) across all logical CPUs

#### Scenario: prometheus/procfs CPUInfo failure does not cause exit 0
- **GIVEN** prometheus/procfs CPUInfo succeeds but prometheus/procfs/sysfs fails to read cpufreq data
- **WHEN** the tool runs
- **THEN** non-frequency fields shall display normally
- **AND** CPU Speed and Max CPU Speed shall display "N/A"
- **AND** the tool shall exit with status 0
