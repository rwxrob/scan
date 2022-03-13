// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package scan implements a non-linear, rune-centric, buffered data,
scanner with user-friendly cursors. The methods of the scanner can be
quickly written by-hand or generated automatically.
*/
package scan

import (
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/rwxrob/structs/qstack"
)

const (

	// EOD is a special value that is returned when the end of data is
	// reached enabling functional parser functions to look for it reliably
	// no matter what is being parsed. Since rune is alias for int32 and
	// Unicode (currently) ends at \U+FFFD we are safe to use the largest
	// possible valid rune value. Passing EOD directly to Expect always
	// stops the scan where it is.
	EOD rune = 1<<31 - 1 // max int32

	// Done means the scanner has reached the end of data and EOD has been
	// set as the last scanned rune.
	Done = 1 << (iota + 1)
)

// R (as in scan.R or "scanner") implements a buffered data, non-linear,
// rune-centric, scanner with cursor and bookmarks for dealing with
// infinite look-ahead/behind designs such as PEG/PEGN. See New for
// creating a usable struct that implements scan.R. The buffer and
// cursor are directly exposed to facilitate higher-performance, direct
// access when needed.
type R struct {

	// Buf is the data buffer providing infinite look-ahead and behind.
	Buf    []byte
	BufLen int

	// Cur is the active current cursor pointing to the Buf data.
	Cur *Cur

	// Snapped contains the latest Cur when Snap was called.
	Snapped *qstack.QStack[*Cur]

	// State allows parser creators to add additional bitwise states as
	// needed. States from 1-999 are reserved but only Done (1) is
	// currently defined. Developers should start their bitwise iotas at
	// 1000.
	State int
}

// New creates a new scan.R instance and initializes it.
func New(i any) (*R, error) {
	s := new(R)
	if err := s.Init(i); err != nil {
		return nil, err
	}
	return s, nil
}

// Init reads all of passed parsable data (io.Reader, string, []byte)
// into buffered memory, scans the first rune, and sets the internals of
// scanner appropriately returning an error if anything happens while
// attempting to read and buffer the data (OOM, etc.).
func (s *R) Init(i any) error {

	s.Cur = new(Cur)
	s.Cur.Pos = Pos{}
	s.Cur.Pos.Line = 1
	s.Cur.Pos.LineRune = 1
	s.Cur.Pos.LineByte = 1
	s.Cur.Pos.Rune = 1

	if err := s.buffer(i); err != nil {
		return err
	}

	r, ln := utf8.DecodeRune(s.Buf) // scan first
	if ln == 0 {
		r = EOD
		s.State |= Done
		return fmt.Errorf("scanner: failed to scan first rune")
	}

	s.Cur.Rune = r
	s.Cur.Len = ln
	s.Cur.Next = ln

	s.Snapped = qstack.New[*Cur]()

	return nil
}

// ---------------------------- marshaling ----------------------------

// String delegates to internal cursor String.
func (s *R) String() string { return s.Cur.String() }

// Print delegates to internal cursor Print.
func (s *R) Print() { s.Cur.Print() }

// Log delegates to internal cursor Log.
func (s *R) Log() { s.Cur.Log() }

// --------------------------------------------------------------------

func (s *R) buffer(i any) error {
	var err error
	switch v := i.(type) {
	case io.Reader:
		s.Buf, err = io.ReadAll(v)
		if err != nil {
			return err
		}
	case string:
		s.Buf = []byte(v)
	case []byte:
		s.Buf = v
	default:
		return fmt.Errorf("scanner: unsupported input type: %T", i)
	}
	s.BufLen = len(s.Buf)
	if s.BufLen == 0 {
		return fmt.Errorf("scanner: no input")
	}
	return err
}

// Scan decodes the next rune and advances the cursor by one.
func (s *R) Scan() {
	if s.Cur.Next == s.BufLen {
		s.Cur.Rune = EOD
		s.State |= Done
		return
	}
	ln := 1
	r := rune(s.Buf[s.Cur.Next])
	if r > utf8.RuneSelf {
		r, ln = utf8.DecodeRune(s.Buf[s.Cur.Next:])
	}
	if ln != 0 {
		s.Cur.Byte = s.Cur.Next
		s.Cur.Pos.LineByte += s.Cur.Len
	} else {
		r = EOD
		s.State |= Done
	}
	s.Cur.Rune = r
	s.Cur.Pos.Rune += 1
	s.Cur.Next += ln
	s.Cur.Pos.LineRune += 1
	s.Cur.Len = ln
}

// Any scans the next n runes advancing n runes forward or returns EOD
// and sets Done state if attempted after already at end of data.
func (s *R) Any(n int) {
	for i := 0; i < n; i++ {
		s.Scan()
	}
}

// Mark returns a copy of the current scanner cursor to preserve like
// a bookmark into the buffer data.
func (s *R) Mark() *Cur {
	if s.Cur == nil {
		return nil
	}
	cp := *s.Cur // force a copy
	return &cp
}

// Jump replaces the internal cursor with a copy of the one passed
// effectively repositioning the scanner's current position in the
// buffered data.
func (s *R) Jump(c *Cur) { nc := *c; s.Cur = &nc }

// Snap pushes a bookmark (as if taken with Mark) onto the Snapped
// stack. Use Back to pop back to the last Snapped.
func (s *R) Snap() {
	cp := *s.Cur // force a copy
	s.Snapped.Push(&cp)
}

// Back pops back to the last Snapped.
func (s *R) Back() {
	if last := s.Snapped.Pop(); last != nil {
		s.Jump(last)
	}
}

// Peek returns a string containing all the runes from the current
// scanner cursor position forward to the number of runes passed.
// If end of data is encountered it will return everything up until that
// point.  Also see Slice and SliceTo.
func (s *R) Peek(n uint) string {
	buf := ""
	pos := s.Cur.Byte
	for c := uint(0); c < n; c++ {
		r, ln := utf8.DecodeRune(s.Buf[pos:])
		if ln == 0 {
			break
		}
		buf += string(r)
		pos += ln
	}
	return buf
}

// PeekSlice returns a string containing all the bytes from the first
// cursor up to the second cursor without changing the cursor position.
func (s *R) PeekSlice(beg *Cur, end *Cur) string {
	return string(s.Buf[beg.Byte:end.Next])
}

// PeekTo returns a string containing all the bytes from the current
// scanner position ahead or behind to the passed cursor position
// without changing positions.
func (s *R) PeekTo(to *Cur) string {
	if to.Byte < s.Cur.Byte {
		return string(s.Buf[to.Byte:s.Cur.Next])
	}
	return string(s.Buf[s.Cur.Byte:to.Next])
}

// NewLine delegates to interval Curs.NewLine to increment the line
// counter to display better parsing status and error information. It is
// up to scanner users to call NewLine explicitly to advance the
// internal cursor when a line is definitively detected.
func (s *R) NewLine() { s.Cur.NewLine() }
