//go:build darwin

package cpu

import "golang.org/x/sys/unix"

// SystemSource is a Source backed by the macOS sysctl interface.
type SystemSource struct{}

func (SystemSource) String(name string) (string, error) { return unix.Sysctl(name) }
func (SystemSource) Uint32(name string) (uint32, error) { return unix.SysctlUint32(name) }
func (SystemSource) Uint64(name string) (uint64, error) { return unix.SysctlUint64(name) }
