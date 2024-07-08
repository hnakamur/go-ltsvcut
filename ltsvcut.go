// Package ltsvcut provides features to cut labels, values from an escaped LTSV
// (Labeled Tab Separated Values) line, and to unescape values.
package ltsvcut

import (
	"bytes"
	"fmt"
)

const labelSeparator = ':'
const fieldSeparator = '\t'

// SkipNFields skips n fields and returns the rest of the input after the nth
// field separator '\t'.
func SkipNFields(input []byte, n int) (rest []byte) {
	for ; n > 0; n-- {
		i := bytes.IndexByte(input, fieldSeparator)
		if i == -1 {
			return nil
		}
		input = input[i+1:]
	}
	return input
}

// CutLabel cuts a label from the beginning of the input and returns the label
// and the rest after the label separator ':'.
func CutLabel(input []byte) (label, rest []byte) {
	i := bytes.IndexByte(input, labelSeparator)
	if i == -1 {
		return nil, input
	}
	return input[:i], input[i+1:]
}

// CutRawValue cuts a raw (escaped) value from the beginning of the input and
// returns the raw value and the rest after the field separator '\t'.
func CutRawValue(input []byte) (rawValue, rest []byte) {
	i := bytes.IndexByte(input, fieldSeparator)
	if i == -1 {
		return input, nil
	}
	return input[:i], input[i+1:]
}

// RawValueForLabel cuts a raw (escaped) value for the label in the input and
// returns the raw value.
//
// It returns nil, false when the label is not found.
func RawValueForLabel(input, label []byte) (rawValue []byte, found bool) {
	rest := input
	for {
		i := bytes.IndexByte(rest, labelSeparator)
		if i == -1 {
			return nil, false
		}

		rawValue = rest[i+1:]
		if bytes.Equal(rest[:i], label) {
			if i := bytes.IndexByte(rawValue, fieldSeparator); i != -1 {
				rawValue = rawValue[:i]
			}
			return rawValue, true
		}

		i = bytes.IndexByte(rawValue, fieldSeparator)
		if i == -1 {
			return nil, false
		}
		rest = rawValue[i+1:]
	}

}

// ValueForLabel cuts and unescape a value for the label in the input and
// returns the unescaped value.
//
// It returns nil, false when the label is not found.
func ValueForLabel(input, label []byte) (value []byte, found bool) {
	rest := input
	for {
		i := bytes.IndexByte(rest, labelSeparator)
		if i == -1 {
			return nil, false
		}

		rawValue := rest[i+1:]
		if bytes.Equal(rest[:i], label) {
			if i := bytes.IndexByte(rawValue, fieldSeparator); i != -1 {
				rawValue = rawValue[:i]
			}
			return UnescapeValue(rawValue), true
		}

		i = bytes.IndexByte(rawValue, fieldSeparator)
		if i == -1 {
			return nil, false
		}
		rest = rawValue[i+1:]
	}
}

// UnescapeValue unescapes a raw value which may be escaped.
// Supported escapes are:
// \t -> tab, \n -> newline, \\ -> backslash.
// It panics when invalid espapes are found.
func UnescapeValue(rawValue []byte) []byte {
	c := countEscape(rawValue)
	if c == 0 {
		return rawValue
	}

	value := make([]byte, 0, len(rawValue)-c)
	i := 0
	for i < len(rawValue) {
		b := rawValue[i]
		i++
		if b == '\\' {
			if i < len(rawValue) {
				switch rawValue[i] {
				case 't':
					b = '\t'
				case 'n':
					b = '\n'
				case '\\':
					b = '\\'
				default:
					panic(fmt.Sprintf("bad escape character in LTSV value: %s", string(rawValue)))
				}
				i++
			} else {
				panic(fmt.Sprintf("no characer after escape in LTSV value: %s", string(rawValue)))
			}
		}
		value = append(value, b)
	}
	return value
}

func countEscape(rawValue []byte) int {
	c := 0
	i := 0
	for i < len(rawValue) {
		if rawValue[i] == '\\' {
			c++
			i += 2
		} else {
			i++
		}
	}
	return c
}
