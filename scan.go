// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package scan provides a rudimentary, performant rune scanner with simple
constructor, trace debugging output, and a package-scope default scanner
instance available for immediate use which can be customized quickly
from init can calling from functions without the need to instantiate an
entirely new one.

Package scan is designed primarily with Parsing Expression Grammars in
mind and assumes the entire bytes buffer to be scanned can be fully
loaded into memory. This dramatically improves scanning speed at the
cost of memory required to scan. Consider the conventional bufio.Scanner
as an alternative when memory pressure is a concern.

*/
package scan

import (
	"fmt"
	"io"
	"log"
	"os"
	"unicode/utf8"
)

// Trace sets the trace for everything that uses this package. Use
// TraceOn/Off for specific scanner tracing.
var Trace int

var DefaultNewLines = []string{"\r\n", "\n"}

// R (to avoid stuttering) implements a buffered data, non-linear,
// rune-centric, scanner with regular expression support
type R struct {
	Buf   []byte   // full buffer for lookahead or behind
	R     rune     // last decoded/scanned rune, maybe >1byte
	B     int      // index pointing beginning of R
	E     int      // index pointing to end (after) R
	NL    []string // all new line variations for Position
	Trace int      // non-zero activates tracing
}

// New is a high-level constructor and alternative to new(scan.R)
// that takes a single optional argument containing any valid Buffer()
// argument. Invalid arguments will fail (not fatal) with log output.
func New(args ...any) *R {
	s := new(R)
	switch len(args) {
	case 2:
		if c, ok := args[1].(Pointer); ok {
			s.Goto(c)
		}
		fallthrough
	case 1:
		s.Buffer(args[0])
	}
	return s
}

// Mark returns a Pointer to a specific rune within a bytes buffer.
// Combine with Goto to jump to positions within the in-memory bytes
// buffer.
func (s *R) Mark() Pointer { return Pointer{&s.Buf, s.R, s.B, s.E} }

// Goto takes a Pointer (usually created with Mark) and positions the
// internal pointer that specific Rune. This can safely be used to go to
// any rune within the buffer provided the pointer is to a valid rune.
// Avoid creating Pointers from scratch preferring Mark since rune
// widths differ significantly.
func (s *R) Goto(c Pointer) { s.R, s.B, s.E = c.R, c.B, c.E }

// CopyEE returns copy (n,m]
func (s *R) CopyEE(m Pointer) string {
	if m.B <= s.B {
		return string(s.Buf[m.E:s.E])
	}
	return string(s.Buf[s.E:m.E])
}

// CopyBB returns copy [n,m]
func (s *R) CopyBE(m Pointer) string {
	if m.B <= s.B {
		return string(s.Buf[m.B:s.E])
	}
	return string(s.Buf[s.B:m.E])
}

// CopyBB returns copy [n,m)
func (s *R) CopyBB(m Pointer) string {
	if m.B <= s.B {
		return string(s.Buf[m.B:s.B])
	}
	return string(s.Buf[s.B:m.B])
}

// CopyEB returns copy (n,m)
func (s *R) CopyEB(m Pointer) string {
	if m.B <= s.B {
		return string(s.Buf[m.E:s.B])
	}
	return string(s.Buf[s.E:m.B])
}

// Open opens the file at path and loads it by passing to Buffer.
func (s *R) Open(path string) error {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return err
	}
	return s.Buffer(f)
}

// Buffer sets the internal bytes buffer (Buf) and resets the existing
// Pointer values to their initial state (null, 0,0). This is useful when
// testing in order to buffer strings as well as content from any
// io.Reader, []byte, []rune, or string.
func (s *R) Buffer(b any) error {
	switch v := b.(type) {
	case string:
		s.Buf = []byte(v)
	case []byte:
		s.Buf = v
	case []rune:
		s.Buf = []byte(string(v))
	case io.Reader:
		b, err := io.ReadAll(v)
		if err != nil {
			return err
		}
		s.Buf = b
	}
	s.R = '\x00'
	s.B = 0
	s.E = 0
	return nil
}

// Print is shorthand for fmt.Println(s).
func (s R) Print() { fmt.Println(s) }

// Log is shorthand for log.Print(s).
func (s R) Log() { log.Println(s) }

// Scan decodes the next rune, setting it to R, and advances position (P)
// by the size of the rune (R) in bytes returning false when there is
// nothing left to scan (similar to bufio.Scanner). Only runes bigger
// than utf8.RuneSelf are decoded since most runes (ASCII) will usually
// be under this number. If s.Trace or scan.Trace is greater than zero
// the scanner trace log output is activated.
//
//
//
func (s *R) Scan() bool {

	if s.E >= len(s.Buf) {
		return false
	}

	ln := 1
	r := rune(s.Buf[s.E])
	if r > utf8.RuneSelf {
		r, ln = utf8.DecodeRune(s.Buf[s.E:])
		if ln == 0 {
			return false
		}
	}

	s.B = s.E
	s.E += ln
	s.R = r

	if s.Trace > 0 || Trace > 0 {
		s.Log()
	}

	return true
}

// Finished returns true if scanner has nothing more to scan.
func (s *R) Finished() bool { return s.E == len(s.Buf) }

// Beginning returns true if and only if the scanner is currently
// pointing to the beginning of the buffer without anything scanned at
// all.
func (s *R) Beginning() bool { return s.E == 0 }

// LookAhead looks for a string match beginning with the first rune
// after the last rune scanned based on the length of the string.
// Returns false if not a match or beyond length of buffer.

// LookBehind looks for a string match ending with the rune before the
// last rune scanned and beginning len(str) before that.  Returns false
// if not a match or beyond beginning of buffer.

// LookAheadExp is same as LookAhead but takes a regular expression.

// LookBehindExp is same sa LookBehind but takes a regular expression.

// Peek returns the specified number of runes starting with the first
// after the last one scanned.

/*
// Peek returns true if the passed string matches from current position
// in the buffer (s.B) forward. Returns false if the string
// would go beyond the length of buffer (len(s.Buf)). Peek does not
// advance the scan.R.
func (s *R) Peek(a string) bool {
	if len(a)+s.E > len(s.Buf) {
		return false
	}
	if string(s.Buf[s.E:s.E+len(a)]) == a {
		return true
	}
	return false
}

// PeekMatch checks for a regular expression match at the current
// position in the buffer (including the last scanned rune (R))
// providing a mechanism for positive and negative lookahead
// expressions. It returns the length of the match.  Successful matches
// might be zero (see regexp.Regexp.FindIndex).  A negative value is
// returned if no match is found. Note that Go regular expressions now
// include the Unicode character classes (ex: \p{L|d}) that should be
// used over dated alternatives (ex: \w).
func (s *R) PeekMatch(re *regexp.Regexp) int {

	loc := re.FindIndex(s.Buf[s.E:])
	if loc == nil {
		return -1
	}

	if loc[0] == 0 {
		return loc[1]
	}

	return -1

}

// Match checks for a regular expression match at the last position in
// the buffer (s.B) providing a mechanism for positive and negative
// lookahead expressions. It returns the length of the match.
// Successful matches might be zero (see regexp.Regexp.FindIndex).
// A negative value is returned if no match is found.  Note that Go
// regular expressions now include the Unicode character classes (ex:
// \p{L|d}) that should be used over dated alternatives (ex: \w).
func (s *R) Match(re *regexp.Regexp) int {
	loc := re.FindIndex(s.Buf[s.B:])
	if loc == nil {
		return -1
	}
	if loc[0] == 0 {
		return loc[1]
	}
	return -1
}

// Pos returns a human-friendly Position for the current location.
// When multiple positions are needed use Positions instead.
func (s R) Pos() Position { return s.Positions(s.E)[0] }

// Positions returns human-friendly Position information (which can easily
// be used to populate a text/template) for each raw byte offset (s.E).
// Only one pass through the buffer (s.Buf) is required to count lines and
// runes since the raw byte position (s.E) is frequently changed
// directly.  Therefore, when multiple positions are wanted, consider
// caching the raw byte positions (s.E) and calling Positions() once for
// all of them.
func (s R) Positions(p ...int) []Position {
	pos := make([]Position, len(p))

	if len(p) == 0 {
		return pos
	}

	if s.NL == nil {
		s.NL = DefaultNewLines
	}

	_rune, line, lbyte, lrune := 1, 1, 1, 1
	_s := R{Buf: s.Buf}
	//_s.Trace++

	for _s.Scan() {

		for _, nl := range s.NL {
			if _s.Is(nl) {
				line++
				_s.E += len(nl) - 1
				_rune += len(nl) - 1
				lbyte = 0
				lrune = 0
				continue
			}
		}

		for i, v := range p {
			if _s.E == v {
				pos[i] = Position{
					Rune:    _s.R,
					BufByte: _s.E,
					BufRune: _rune,
					Line:    line,
					LByte:   lbyte,
					LRune:   lrune,
				}
			}
		}

		rlen := len([]byte(string(s.R)))
		lbyte += rlen
		lrune++
		_rune++

	}

	return pos
}
*/
