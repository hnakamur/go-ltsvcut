package ltsvcut_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/hnakamur/go-ltsvcut"
)

func TestCutLabelAndCutRawValueAndSkipNFields(t *testing.T) {
	input := []byte("time:2024-01-30T15:23:46.123Z\treq:GET / HTTP/1.1\tstatus:200\tua:name with escaped \\ttab, \\\\backslash, and \\nnewline")
	label, rest := ltsvcut.CutLabel(input)
	if got, want := label, []byte("time"); !bytes.Equal(got, want) {
		t.Errorf("first label mismatched, got=%s, want=%s", string(got), string(want))
	}
	rawValue, rest := ltsvcut.CutRawValue(rest)
	if got, want := ltsvcut.UnescapeValue(rawValue), []byte("2024-01-30T15:23:46.123Z"); !bytes.Equal(got, want) {
		t.Errorf("first value mismatched, got=%s, want=%s", string(got), string(want))
	}
	rest = ltsvcut.SkipNFields(rest, 2)
	label, rest = ltsvcut.CutLabel(rest)
	if got, want := label, []byte("ua"); !bytes.Equal(got, want) {
		t.Errorf("second label mismatched, got=%s, want=%s", string(got), string(want))
	}
	rawValue, _ = ltsvcut.CutRawValue(rest)
	if got, want := ltsvcut.UnescapeValue(rawValue), []byte("name with escaped \ttab, \\backslash, and \nnewline"); !bytes.Equal(got, want) {
		t.Errorf("second value mismatched, got=%s, want=%s", string(got), string(want))
	}
}

func TestRawValueForLabel(t *testing.T) {
	input := []byte("time:2024-01-30T15:23:46.123Z\treq:GET / HTTP/1.1\tstatus:200\tua:name with escaped \\ttab, \\\\backslash, and \\nnewline")

	rawValue, found := ltsvcut.RawValueForLabel(input, []byte("ua"))
	if got, want := found, true; got != want {
		t.Errorf("found mismatch, got=%v, want=%v", got, want)
	}
	if got, want := rawValue, []byte("name with escaped \\ttab, \\\\backslash, and \\nnewline"); !bytes.Equal(got, want) {
		t.Errorf("value mismatch, got=%s, want=%s", string(got), string(want))
	}

	_, found = ltsvcut.RawValueForLabel(input, []byte("no_such_label"))
	if got, want := found, false; got != want {
		t.Errorf("found mismatch, got=%v, want=%v", got, want)
	}
}

func TestValueForLabel(t *testing.T) {
	input := []byte("time:2024-01-30T15:23:46.123Z\treq:GET / HTTP/1.1\tstatus:200\tua:name with escaped \\ttab, \\\\backslash, and \\nnewline")

	value, found := ltsvcut.ValueForLabel(input, []byte("ua"))
	if got, want := found, true; got != want {
		t.Errorf("found mismatch, got=%v, want=%v", got, want)
	}
	if got, want := value, []byte("name with escaped \ttab, \\backslash, and \nnewline"); !bytes.Equal(got, want) {
		t.Errorf("value mismatch, got=%s, want=%s", string(got), string(want))
	}

	_, found = ltsvcut.ValueForLabel(input, []byte("no_such_label"))
	if got, want := found, false; got != want {
		t.Errorf("found mismatch, got=%v, want=%v", got, want)
	}
}

func Example() {
	input := []byte("time:2024-01-30T15:23:46.123Z\tua:name with escaped\\ttab,\\\\backslash, and\\nnewline\n" +
		"time:2024-01-30T15:23:46.456Z\tua:my agent\n")
	r := bufio.NewScanner(bytes.NewReader(input))
	for r.Scan() {
		rest := r.Bytes()
		for {
			var label, rawValue []byte
			label, rest = ltsvcut.CutLabel(rest)
			if label == nil {
				break
			}
			rawValue, rest = ltsvcut.CutRawValue(rest)
			value := ltsvcut.UnescapeValue(rawValue)
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
