package cpu

import (
	"errors"
	"testing"
)

type fakeSource struct {
	strings map[string]string
	uint32s map[string]uint32
	uint64s map[string]uint64
}

var errUnknownOID = errors.New("unknown oid")

func (f fakeSource) String(name string) (string, error) {
	if v, ok := f.strings[name]; ok {
		return v, nil
	}
	return "", errUnknownOID
}

func (f fakeSource) Uint32(name string) (uint32, error) {
	if v, ok := f.uint32s[name]; ok {
		return v, nil
	}
	return 0, errUnknownOID
}

func (f fakeSource) Uint64(name string) (uint64, error) {
	if v, ok := f.uint64s[name]; ok {
		return v, nil
	}
	return 0, errUnknownOID
}

func TestVendorFromModel(t *testing.T) {
	tests := []struct {
		name  string
		model string
		want  string
	}{
		{"apple silicon", "Apple M5 Pro", "Apple"},
		{"intel with trademarks", "Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz", "Intel"},
		{"leading spaces", "   Intel(R) Xeon(R) CPU", "Intel"},
		{"empty", "", ""},
		{"only spaces", "   ", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vendorFromModel(tt.model); got != tt.want {
				t.Errorf("vendorFromModel(%q) = %q, want %q", tt.model, got, tt.want)
			}
		})
	}
}

func TestRead(t *testing.T) {
	tests := []struct {
		name string
		src  fakeSource
		want CPU
	}{
		{
			name: "apple silicon, 18 physical and 18 logical, no frequency",
			src: fakeSource{
				strings: map[string]string{"machdep.cpu.brand_string": "Apple M5 Pro"},
				uint32s: map[string]uint32{"hw.physicalcpu": 18, "hw.logicalcpu": 18},
			},
			want: CPU{Vendor: "Apple", Model: "Apple M5 Pro", Physical: 18, Logical: 18},
		},
		{
			name: "intel, 6 physical and 12 logical, frequency present",
			src: fakeSource{
				strings: map[string]string{"machdep.cpu.brand_string": "Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz"},
				uint32s: map[string]uint32{"hw.physicalcpu": 6, "hw.logicalcpu": 12},
				uint64s: map[string]uint64{"hw.cpufrequency_max": 2600000000},
			},
			want: CPU{
				Vendor:   "Intel",
				Model:    "Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz",
				SpeedHz:  2600000000,
				Physical: 6,
				Logical:  12,
			},
		},
		{
			name: "empty brand string",
			src: fakeSource{
				strings: map[string]string{"machdep.cpu.brand_string": ""},
				uint32s: map[string]uint32{"hw.physicalcpu": 4, "hw.logicalcpu": 8},
			},
			want: CPU{Physical: 4, Logical: 8},
		},
		{
			name: "nothing available",
			src:  fakeSource{},
			want: CPU{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Read(tt.src); got != tt.want {
				t.Errorf("Read() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
