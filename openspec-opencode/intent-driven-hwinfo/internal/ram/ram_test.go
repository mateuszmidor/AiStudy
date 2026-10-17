package ram

import (
	"errors"
	"testing"
)

func TestFormatTotalLineConvertsBytesToGB(t *testing.T) {
	got := FormatTotalLine(68719476736, true)
	want := "Total Memory: 64 GB"
	if got != want {
		t.Errorf("FormatTotalLine(68719476736, true) = %q, want %q", got, want)
	}
}

func TestFormatTotalLineTruncatesTowardZero(t *testing.T) {
	// One byte less than 16 GiB: the decimal portion is dropped, never rounded up.
	got := FormatTotalLine(17179869183, true)
	want := "Total Memory: 15 GB"
	if got != want {
		t.Errorf("FormatTotalLine(17179869183, true) = %q, want %q", got, want)
	}
}

func TestFormatTotalLineWhenUnavailable(t *testing.T) {
	for _, total := range []uint64{0, 68719476736} {
		got := FormatTotalLine(total, false)
		want := "Total Memory: N/A"
		if got != want {
			t.Errorf("FormatTotalLine(%d, false) = %q, want %q", total, got, want)
		}
	}
}

func TestTotalPassesReadBytesThrough(t *testing.T) {
	original := readTotal
	defer func() { readTotal = original }()
	readTotal = func() (uint64, error) { return 68719476736, nil }

	got, err := Total()
	if err != nil {
		t.Fatalf("Total() returned unexpected error: %v", err)
	}
	if got != 68719476736 {
		t.Errorf("Total() = %d, want %d", got, 68719476736)
	}
}

func TestTotalPropagatesReadError(t *testing.T) {
	original := readTotal
	defer func() { readTotal = original }()
	wantErr := errors.New("read failed")
	readTotal = func() (uint64, error) { return 0, wantErr }

	if _, err := Total(); !errors.Is(err, wantErr) {
		t.Errorf("Total() error = %v, want %v", err, wantErr)
	}
}
