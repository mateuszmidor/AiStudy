package main

import (
	"fmt"

	"hwinfo/internal/cpu"
)

func main() {
	fmt.Print(cpu.Format(cpu.Read(cpu.SystemSource{})))
}
