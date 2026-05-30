package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/prometheus/procfs/sysfs"
	"github.com/shirou/gopsutil/v4/cpu"
)

func main() {
	info, err := cpu.Info()
	if err != nil || len(info) == 0 {
		fmt.Fprintf(os.Stderr, "Error: unable to read CPU information\n")
		os.Exit(1)
	}

	manufacturer := info[0].VendorID
	model := info[0].ModelName
	architecture := runtime.GOARCH

	physicalCores, _ := cpu.Counts(false)
	logicalCores, _ := cpu.Counts(true)

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
}
