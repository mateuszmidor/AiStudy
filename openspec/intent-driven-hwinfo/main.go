package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/sysfs"
	"github.com/shirou/gopsutil/v4/sensors"
)

func cpuTemperature() string {
	temps, err := sensors.SensorsTemperatures()
	if err != nil {
		return "N/A"
	}

	var total, count float64
	for _, t := range temps {
		key := strings.ToLower(t.SensorKey)
		if strings.Contains(key, "k10temp") ||
			strings.Contains(key, "coretemp") ||
			strings.Contains(key, "zen") {
			total += t.Temperature
			count++
		}
	}

	if count == 0 {
		return "N/A"
	}

	return fmt.Sprintf("%d", int(total/count))
}

func main() {
	pfs, err := procfs.NewDefaultFS()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: unable to read CPU information\n")
		os.Exit(1)
	}
	cpuInfo, err := pfs.CPUInfo()
	if err != nil || len(cpuInfo) == 0 {
		fmt.Fprintf(os.Stderr, "Error: unable to read CPU information\n")
		os.Exit(1)
	}

	manufacturer := cpuInfo[0].VendorID
	model := cpuInfo[0].ModelName
	architecture := runtime.GOARCH

	seen := make(map[string]struct{})
	for _, c := range cpuInfo {
		seen[c.PhysicalID+":"+c.CoreID] = struct{}{}
	}
	physicalCores := len(seen)
	logicalCores := len(cpuInfo)

	avgMhzStr := "N/A"
	maxMhzStr := "N/A"

	fs, err := sysfs.NewDefaultFS()
	if err == nil {
		cpuFreqs, err := fs.SystemCpufreq()
		if err == nil {
			var currentTotal, maxTotal float64
			var currentCount, maxCount int
			for _, cf := range cpuFreqs {
				if cf.ScalingCurrentFrequency != nil {
					currentTotal += float64(*cf.ScalingCurrentFrequency) / 1000.0
					currentCount++
				}
				if cf.CpuinfoMaximumFrequency != nil {
					maxTotal += float64(*cf.CpuinfoMaximumFrequency) / 1000.0
					maxCount++
				}
			}
			if currentCount > 0 {
				avgMhzStr = fmt.Sprintf("%d", int(currentTotal/float64(currentCount)))
			}
			if maxCount > 0 {
				maxMhzStr = fmt.Sprintf("%d", int(maxTotal/float64(maxCount)))
			}
		}
	}

	if manufacturer == "" {
		manufacturer = "N/A"
	}
	if model == "" {
		model = "N/A"
	}
	if architecture == "" {
		architecture = "N/A"
	}

	fmt.Println("HWInfo - CPU Information")
	fmt.Println("------------------------")
	fmt.Printf("CPU Manufacturer: %s\n", manufacturer)
	fmt.Printf("CPU Model: %s\n", model)
	fmt.Printf("CPU Architecture: %s\n", architecture)

	if physicalCores > 0 {
		fmt.Printf("Physical Cores: %d\n", physicalCores)
	} else {
		fmt.Println("Physical Cores: N/A")
	}

	if logicalCores > 0 {
		fmt.Printf("Logical Cores: %d\n", logicalCores)
	} else {
		fmt.Println("Logical Cores: N/A")
	}

	fmt.Printf("CPU Speed (MHz): %s\n", avgMhzStr)
	fmt.Printf("Max CPU Speed (MHz): %s\n", maxMhzStr)

	cpuTemp := cpuTemperature()
	if cpuTemp == "N/A" {
		fmt.Println("CPU Temperature (Avg): N/A")
	} else {
		fmt.Printf("CPU Temperature (Avg): %s°C\n", cpuTemp)
	}
}
