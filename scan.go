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
	"log"
	"regexp"
	"unicode/utf8"
)

// Trace activates tracing for anything using the package. This is
// sometimes more convenient when an application uses the package but
// does not give access to the equivalent R.Trace property.
var Trace int

// ViewLen sets the number of bytes to view before eliding the rest.
var ViewLen = 20

// R (as in scan.R or "scanner") implements a buffered data, non-linear,
// rune-centric, scanner with regular expression support. Keep in mind
// that if and when you change the position (P) directly that rune (R) will not
// itself be updated as it is only updated by calling Scan. Often an
// update to the rune (R) as well would be inconsequential, even wasteful.
type R struct {
	B      []byte // full buffer for lookahead or behind
	P      int    // current position in the buffer
	R      rune   // updated by Scan
	Trace  int
	Errors []error
}

// String implements fmt.Stringer with simply the position (P) and
// quoted rune (R) along with its Unicode.
func (s R) String() string {
	end := s.P + ViewLen
	elided := "..."
	if end > len(s.B) {
		end = len(s.B)
		elided = ""
	}
	return fmt.Sprintf("%v %q %q%v",
		s.P, s.R, s.B[s.P:end], elided)
}

// Print is shorthand for fmt.Println(s).
func (s R) Print() { fmt.Println(s) }

// Log is shorthand for log.Print(s).
func (s R) Log() { log.Println(s) }

// Scan decodes the next rune, setting it to R, and advances position
// (P) by the size of the rune (R) in bytes returning false then there
// is nothing left to scan. Only runes bigger than utf8.RuneSelf are
// decoded since most runes (ASCII) will usually be under this number.
func (s *R) Scan() bool {

	if s.P >= len(s.B) {
		return false
	}

	ln := 1
	r := rune(s.B[s.P])
	if r > utf8.RuneSelf {
		r, ln = utf8.DecodeRune(s.B[s.P:])
		if ln == 0 {
			return false
		}
	}

	s.P += ln
	s.R = r

	if s.Trace > 0 || Trace > 0 {
		s.Log()
	}

	return true
}

// Peek returns true if the passed string matches the current position
// including the last s.R (s.P-1) in the buffer. Returns false if
// the string would go beyond the length of buffer (len(s.B)).
func (s *R) Peek(a string) bool {
	if len(a)+s.P-1 > len(s.B) {
		return false
	}
	if string(s.B[s.P-1:s.P-1+len(a)]) == a {
		return true
	}
	return false
}

// Match checks for a regular expression match at the current position
// in the buffer providing a mechanism for positive and negative
// lookahead expressions (which includes the current s.R, s.P-1).  It
// returns the length of the match. Successful matches might be zero
// (see regexp.Regexp.FindIndex). A negative value is returned if no
// match is found. Keep in mind that Note that Go regular expressions
// now include the Unicode character classes (ex: \p{L}) that should be
// used over dated alternatives (ex: \w).
func (s *R) Match(re *regexp.Regexp) int {
	loc := re.FindIndex(s.B[s.P-1:])
	if loc == nil {
		return -1
	}
	if loc[0] == 0 {
		return loc[1]
	}
	return -1
}
