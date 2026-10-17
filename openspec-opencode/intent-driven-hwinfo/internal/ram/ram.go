// Package ram reports the machine's total physical memory.
package ram

import (
	"strconv"

	"github.com/shirou/gopsutil/v4/mem"
)

// readTotal is the seam that reads total physical memory; tests replace it to
// supply bytes without touching the OS.
var readTotal = func() (uint64, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vm.Total, nil
}

// Total reports the total physical memory of the machine in bytes.
func Total() (uint64, error) {
	return readTotal()
}

// FormatTotalLine renders the output line for a total, or the unavailable form
// when available is false. The byte count is divided by 1024^3 and truncated
// toward zero, so it never rounds up.
func FormatTotalLine(total uint64, available bool) string {
	const label = "Total Memory: "
	if !available {
		return label + "N/A"
	}
	return label + strconv.FormatUint(total/(1024*1024*1024), 10) + " GB"
}
