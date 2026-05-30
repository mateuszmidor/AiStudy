package main

import (
	"fmt"
	"os"
	"runtime"

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

	var totalMhz float64
	var count int
	for _, c := range info {
		if c.Mhz > 0 {
			totalMhz += c.Mhz
			count++
		}
	}

	var avgMhzStr string
	if count > 0 {
		avgMhz := totalMhz / float64(count)
		avgMhzStr = fmt.Sprintf("%v", avgMhz)
	} else {
		avgMhzStr = "N/A"
	}

	maxMhzStr := "N/A"

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
