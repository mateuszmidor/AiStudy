package cpu

import (
	"fmt"
	"strings"
)

const unavailable = "unavailable"

// Format renders c as a labeled text section.
func Format(c CPU) string {
	var b strings.Builder
	b.WriteString("CPU\n")
	fmt.Fprintf(&b, "  Vendor: %s\n", orUnavailable(c.Vendor))
	fmt.Fprintf(&b, "  Model: %s\n", orUnavailable(c.Model))
	fmt.Fprintf(&b, "  Speed: %s\n", formatSpeed(c.SpeedHz))
	fmt.Fprintf(&b, "  Cores: %s\n", formatCores(c.Physical, c.Logical))
	return b.String()
}

func orUnavailable(s string) string {
	if s == "" {
		return unavailable
	}
	return s
}

func formatSpeed(hz uint64) string {
	if hz == 0 {
		return unavailable
	}
	return fmt.Sprintf("%.2f GHz", float64(hz)/1e9)
}

func formatCores(physical, logical int) string {
	if physical == 0 || logical == 0 {
		return unavailable
	}
	return fmt.Sprintf("%d physical, %d logical", physical, logical)
}
