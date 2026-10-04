// Package cpu reads and formats basic CPU information.
package cpu

import "strings"

// Source looks up system values by name (for example sysctl keys).
// A lookup returns an error when the value does not exist on this machine.
type Source interface {
	String(name string) (string, error)
	Uint32(name string) (uint32, error)
	Uint64(name string) (uint64, error)
}

// CPU holds basic CPU facts. A zero value means the fact is unavailable.
type CPU struct {
	Vendor   string
	Model    string
	SpeedHz  uint64
	Physical int
	Logical  int
}

var trademarks = strings.NewReplacer("(R)", "", "(TM)", "", "(tm)", "")

// Read collects CPU facts from src. Lookups that fail leave the field unavailable.
func Read(src Source) CPU {
	var c CPU

	if brand, err := src.String("machdep.cpu.brand_string"); err == nil {
		c.Model = strings.TrimSpace(brand)
		c.Vendor = vendorFromModel(c.Model)
	}
	if n, err := src.Uint32("hw.physicalcpu"); err == nil {
		c.Physical = int(n)
	}
	if n, err := src.Uint32("hw.logicalcpu"); err == nil {
		c.Logical = int(n)
	}
	if hz, err := src.Uint64("hw.cpufrequency_max"); err == nil {
		c.SpeedHz = hz
	}

	return c
}

// vendorFromModel returns the first word of the model with trademark marks
// removed, or "" when there is none.
func vendorFromModel(model string) string {
	words := strings.Fields(model)
	if len(words) == 0 {
		return ""
	}
	return trademarks.Replace(words[0])
}
