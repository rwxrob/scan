// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package scan implements a non-linear, rune-centric, buffered data,
scanner with strong support for regular expressions. The methods of the
scanner can be quickly written by-hand or generated automatically.
*/
package scan

import (
	"fmt"
	"unicode/utf8"
)

// R (as in scan.R or "scanner") implements a buffered data, non-linear,
// rune-centric, scanner with regular expression support.
type R struct {
	Buf  []byte
	Pos  int
	Rune rune
}

// String implements fmt.Stringer with simply the Pos and quoted Rune
// along with its Unicode.
func (s *R) String() string {
	return fmt.Sprintf("%v %U %q", s.Pos, s.Rune, s.Rune)
}

// Print is shorthand for fmt.Println(s).
func (s *R) Print() { fmt.Println(s) }

// Scan decodes the next rune, setting it to Rune, and advances Pos by
// the size of the Rune in bytes returning false then there is nothing
// left to scan. Only runes bigger than utf8.RuneSelf are decoded since
// most runes (ASCII) will usually be under this number.
func (s *R) Scan() bool {
	if len(s.Buf) == s.Pos {
		return false
	}
	ln := 1
	r := rune(s.Buf[s.Pos])
	if r > utf8.RuneSelf {
		r, ln = utf8.DecodeRune(s.Buf[s.Pos:])
		if ln == 0 {
			return false
		}
	}
	s.Pos += ln
	s.Rune = r
	return true
}

// ScanN calls Scan n number of times stopping and returning false if
// and when Scan returns false.
func (s *R) ScanN(n int) bool {
	for i := 0; i < n; i++ {
		if !s.Scan() {
			return false
		}
	}
	return true
}

// Is returns true if the passed string matches the current position in
// the buffer.
func (s *R) Is(a string) bool {
	if len(a)+s.Pos > len(s.Buf) {
		return false
	}
	if string(s.Buf[s.Pos:s.Pos+len(a)]) == a {
		return true
	}
	return false
}
