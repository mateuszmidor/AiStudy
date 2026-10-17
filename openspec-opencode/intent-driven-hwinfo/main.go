// Command hwinfo prints the machine's total physical memory as one line on
// stdout. It is never fatal: a failed read and a zero total both print the
// unavailable form and exit 0.
package main

import (
	"fmt"

	"hwinfo/internal/ram"
)

func main() {
	total, err := ram.Total()
	// A read failure and a successful read of 0 bytes are distinct conditions
	// that both mean the total is unavailable.
	available := err == nil && total > 0
	fmt.Println(ram.FormatTotalLine(total, available))
}
