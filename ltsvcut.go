// Package ltsvcut provides features to cut labels, values from an escaped LTSV
// (Labeled Tab Separated Values) line, and to unescape values.
package ltsvcut

import (
	"bytes"
	"fmt"
)

// Cutter is a struct which holds a line and the current position.
// It can be reused for another line.
type Cutter struct {
	line []byte
	pos  int
}

// SetLine set the current line and reset the position to 0.
func (c *Cutter) SetLine(line []byte) {
	c.line = line
	c.pos = 0
}

// NextLabel returns the next label or nil if no label found after the current
// position.
func (c *Cutter) NextLabel() []byte {
	i := bytes.IndexByte(c.line[c.pos:], ':')
	if i == -1 {
		return nil
	}
	label := c.line[c.pos : c.pos+i]
	c.pos += i + 1
	return label
}

// UnescapedValue returns the unescaped value after the current position.
func (c *Cutter) UnescapedValue() []byte {
	return UnescapeValue(c.RawValue())
}

// UnescapedValue returns the raw value (which may be escaped) after the current position.
func (c *Cutter) RawValue() []byte {
	i := bytes.IndexByte(c.line[c.pos:], '\t')
	var rawValue []byte
	if i == -1 {
		rawValue = c.line[c.pos:]
		c.pos = len(c.line)
	} else {
		rawValue = c.line[c.pos : c.pos+i]
		c.pos += i + 1
	}
	return rawValue
}

// UnescapeValue unescaped a raw value which may be escaped.
// Supported escapes are:
// \t -> tab, \n -> newline, \\ -> backslash.
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
