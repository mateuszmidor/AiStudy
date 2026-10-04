package cpu

import (
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name string
		cpu  CPU
		want []string
	}{
		{
			name: "apple silicon without speed",
			cpu:  CPU{Vendor: "Apple", Model: "Apple M5 Pro", Physical: 18, Logical: 18},
			want: []string{
				"Vendor: Apple",
				"Model: Apple M5 Pro",
				"Speed: unavailable",
				"Cores: 18 physical, 18 logical",
			},
		},
		{
			name: "more logical than physical with speed",
			cpu:  CPU{Vendor: "Intel", Model: "Intel(R) Core(TM) i7", SpeedHz: 2600000000, Physical: 6, Logical: 12},
			want: []string{
				"Vendor: Intel",
				"Speed: 2.60 GHz",
				"Cores: 6 physical, 12 logical",
			},
		},
		{
			name: "nothing available",
			cpu:  CPU{},
			want: []string{
				"Vendor: unavailable",
				"Model: unavailable",
				"Speed: unavailable",
				"Cores: unavailable",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Format(tt.cpu)
			for _, line := range tt.want {
				if !strings.Contains(got, line+"\n") {
					t.Errorf("Format() missing line %q in:\n%s", line, got)
				}
			}
		})
	}
}

func TestFormatSpeed(t *testing.T) {
	tests := []struct {
		hz   uint64
		want string
	}{
		{2600000000, "2.60 GHz"},
		{3490000000, "3.49 GHz"},
		{0, "unavailable"},
	}
	for _, tt := range tests {
		if got := formatSpeed(tt.hz); got != tt.want {
			t.Errorf("formatSpeed(%d) = %q, want %q", tt.hz, got, tt.want)
		}
	}
}
