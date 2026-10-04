# CPU Temperature

## Purpose

Read CPU die/package temperature from hardware sensors and display it as an average value in degrees Celsius. Sourced via the gopsutil/v4/sensors Go package.

## ADDED Requirements

### Requirement: Report average CPU temperature
The tool SHALL read CPU temperature from hardware sensors and display the average in degrees Celsius.

#### Scenario: Average CPU temperature displayed successfully
- **GIVEN** the system has a CPU temperature sensor (e.g., k10temp, coretemp) reported by gopsutil
- **WHEN** the tool reads CPU temperature via gopsutil/v4/sensors
- **THEN** it shall filter sensor results for CPU sensors (matching "k10temp", "coretemp", or "zen")
- **AND** it shall compute the arithmetic mean of all matched sensor values
- **AND** it shall display the value as integer degrees Celsius (decimal portion truncated)
- **AND** it shall format the field as "CPU Temperature (Avg): <value>°C"

#### Scenario: CPU sensor unavailable shows N/A
- **GIVEN** no CPU temperature sensor is reported by gopsutil (e.g., virtual machine, missing kernel driver)
- **WHEN** the tool reads CPU temperature
- **THEN** the field shall display "CPU Temperature (Avg): N/A"
- **AND** the tool shall exit with status 0

#### Scenario: gopsutil error shows N/A
- **GIVEN** gopsutil/v4/sensors returns an error (e.g., permission denied, unsupported platform)
- **WHEN** the tool reads CPU temperature
- **THEN** the field shall display "CPU Temperature (Avg): N/A"
- **AND** the tool shall exit with status 0

### Requirement: CPU Temperature data source
The tool SHALL read CPU temperature using the gopsutil/v4/sensors Go package, filtering for CPU-specific sensors by key name.

#### Scenario: Match k10temp for AMD CPUs
- **GIVEN** the system has an AMD CPU with the k10temp driver loaded
- **WHEN** the tool fetches temperatures via gopsutil
- **AND** it finds a sensor key containing "k10temp"
- **THEN** it shall include that sensor's value in the average

#### Scenario: Match coretemp for Intel CPUs
- **GIVEN** the system has an Intel CPU with the coretemp driver loaded
- **WHEN** the tool fetches temperatures via gopsutil
- **AND** it finds a sensor key containing "coretemp"
- **THEN** it shall include that sensor's value in the average
