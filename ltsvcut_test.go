package ltsvcut_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/hnakamur/ltsvcut"
)

func TestCutter(t *testing.T) {
	input := []byte("time:2024-01-30T15:23:46.123Z\tua:name with escaped \\ttab, \\\\backslash, and \\nnewline")
	var cutter ltsvcut.Cutter
	cutter.SetLine(input)
	if got, want := cutter.NextLabel(), []byte("time"); !bytes.Equal(got, want) {
		t.Errorf("first label mismatched, got=%s, want=%s", string(got), string(want))
	}
	if got, want := cutter.UnescapedValue(), []byte("2024-01-30T15:23:46.123Z"); !bytes.Equal(got, want) {
		t.Errorf("first value mismatched, got=%s, want=%s", string(got), string(want))
	}
	if got, want := cutter.NextLabel(), []byte("ua"); !bytes.Equal(got, want) {
		t.Errorf("second label mismatched, got=%s, want=%s", string(got), string(want))
	}
	if got, want := cutter.UnescapedValue(), []byte("name with escaped \ttab, \\backslash, and \nnewline"); !bytes.Equal(got, want) {
		t.Errorf("second value mismatched, got=%s, want=%s", string(got), string(want))
	}
}

func ExampleCutter() {
	input := []byte("time:2024-01-30T15:23:46.123Z\tua:name with escaped\\ttab,\\\\backslash, and\\nnewline\n" +
		"time:2024-01-30T15:23:46.456Z\tua:my agent\n")
	r := bufio.NewScanner(bytes.NewReader(input))
	var cutter ltsvcut.Cutter
	for r.Scan() {
		cutter.SetLine(r.Bytes())
		for {
			label := cutter.NextLabel()
			if label == nil {
				break
			}
			value := cutter.UnescapedValue()
			fmt.Printf("label=%s, value=%s\n", label, value)
		}
		fmt.Printf("---\n")
	}
	if err := r.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner err: %s", err)
	}
	// Output:
	// label=time, value=2024-01-30T15:23:46.123Z
	// label=ua, value=name with escaped	tab,\backslash, and
	// newline
	// ---
	// label=time, value=2024-01-30T15:23:46.456Z
	// label=ua, value=my agent
	// ---
}
