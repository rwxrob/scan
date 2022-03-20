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
	"log"
	"unicode/utf8"

	z "github.com/rwxrob/scan/is"
	"github.com/rwxrob/scan/tk"
	"github.com/rwxrob/structs/qstack"
	"github.com/rwxrob/structs/tree"
)

const (
	EOD = 1 << (iota + 1) // end of data has been reached
)

// ---------------------------- scan.Error ----------------------------

// Error captures an error at a specific location.
type Error struct {
	Msg  string
	At   *Cur
	What any
}

// String fulfills the fmt.Stringer interface.
func (e *Error) String() string {
	return fmt.Sprintf(`%v at %v`, e.Msg, e.At)
}

// ------------------------------ scan.R ------------------------------

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

	// Last contains the last cursor position before Scan.
	Last *Cur

	// Snapped contains the latest Cur when Snap was called.
	Snapped *qstack.QS[*Cur]

	// State allows parser creators to add additional bitwise states as
	// needed. States from 1-999 are reserved. Developers should start
	// their bitwise iotas at 1000.
	State int

	// Err contains errors encountered while scanning expressions.
	Err *qstack.QS[*Error]

	// Tree contains the rooted node tree created by z.P scan.X expressions.
	Tree *tree.E[string]

	// Nodes is used to construct what ultimately becomes the Tree.Root
	// after all z.P scan.X parsing completes. Init pushes Tree.Root onto
	// it to begin. It is public so that first-class function scan.X
	// expression authors can make use of it for their own parsing
	// possibilities.
	Nodes *qstack.QS[*tree.Node[string]]

	tracex bool
}

// New creates a new scan.R instance and initializes it pushing an error
// to Err if any are encountered.
func New(i any) *R {
	s := new(R)
	s.Init(i)
	return s
}

// Init reads all of passed parsable data (io.Reader, string, []byte)
// into buffered memory, scans the first rune, and sets the internals of
// scanner appropriately pushing an error to Err if anything happens
// while attempting to read and buffer the data (OOM, etc.). Init sets
// Snapped, Err, Tree, Nodes, and Cur to their initialized values.
func (s *R) Init(i any) {

	s.Snapped = qstack.New[*Cur]()
	s.Err = qstack.New[*Error]()
	s.Tree = tree.New[string]()
	s.Nodes = qstack.New[*tree.Node[string]]()
	s.Nodes.Push(s.Tree.Root)
	s.tracex = false

	s.Cur = new(Cur)
	s.Cur.Pos = Pos{}
	s.Cur.Pos.Line = 1
	s.Cur.Pos.LineRune = 1
	s.Cur.Pos.LineByte = 1
	s.Cur.Pos.Rune = 1

	s.buffer(i)
	if s.Err.Len > 0 {
		return
	}

	r, ln := utf8.DecodeRune(s.Buf)
	if ln == 0 {
		r = tk.EOD
		s.State |= EOD
		s.Errorf(nil, "init: failed to scan first rune")
		return
	}

	s.Cur.Rune = r
	s.Cur.Len = ln
	s.Cur.Next = ln

}

// TraceX activates scan.X interpreter tracing providing crucial
// visibility for complex expression debugging. Scan.X expressions are
// designed for rapid development of complex grammars and expressions.
// Later, when and if needed, those expressions that can be used to
// generate highly optimized parsers that do not depend on the scan.X
// interpreter.
func (s *R) TraceX() { s.tracex = true }

// NoTraceX disables tracing of scan.X expressions. See TraceX.
func (s *R) NoTraceX() { s.tracex = false }

// Error pushes the message from the passed error onto the Err stack.
func (s *R) Error(err error) {
	s.Err.Push(&Error{fmt.Sprintf(`%v`, err), s.Mark(), nil})
}

// Errorf pushes a new formatted error on Err stack. The last error
// is always displayed with the scan.R is marshaled/printed as a string.
// The first argument may be the context (item) that caused the error.
func (s *R) Errorf(t any, tpl string, i ...any) {
	msg := fmt.Sprintf(tpl, i...)
	s.Err.Push(&Error{msg, s.Mark(), t})
}

func (s *R) ClearLastError() {
	err := s.Err.Pop()
	if _, is := err.What.(z.P); is {
		n := s.Nodes.Pop()
		n.Cut()
	}
}

// String prints the last error and position.
func (s *R) String() string {
	if s.Err.Len > 0 {
		return s.Err.Pop().String()
	}
	return s.Cur.String()
}

// Print delegates to internal cursor Print.
func (s *R) Print() { fmt.Println(s.String()) }

// Log delegates to internal cursor Log.
func (s *R) Log() { log.Print(s.String()) }

func (s *R) buffer(i any) {
	switch v := i.(type) {
	case io.Reader:
		var err error
		s.Buf, err = io.ReadAll(v)
		if err != nil {
			s.Error(err)
			return
		}
	case string:
		s.Buf = []byte(v)
	case []byte:
		s.Buf = v
	default:
		s.Errorf(v, "buffer: unsupported type: %T", i)
		return
	}
	s.BufLen = len(s.Buf)
	if s.BufLen == 0 {
		s.Errorf(nil, "buffer: no input")
		return
	}
}

// Scan decodes the next rune and advances the cursor by one.  If the
// scan exceeds BufLen then Cur.Rune is set to tk.EOD, EOD State
// is set, and Scan returns false.
func (s *R) Scan() bool {
	s.Last = s.Mark()
	if s.Cur.Next == s.BufLen {
		s.Cur.Rune = tk.EOD
		s.State |= EOD
		return false
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
		r = tk.EOD
		s.State |= EOD
	}
	s.Cur.Rune = r
	s.Cur.Pos.Rune += 1
	s.Cur.Next += ln
	s.Cur.Pos.LineRune += 1
	s.Cur.Len = ln
	if r == tk.EOD {
		return false
	}
	return true
}

// Any calls Scan n number of times stopping if end of data reached.
func (s *R) Any(n int) bool {
	for i := 0; i < n; i++ {
		if !s.Scan() {
			return false
		}
	}
	return true
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
func (s *R) Jump(c *Cur) {
	if c == nil {
		return
	}
	nc := *c
	s.Cur = &nc
}

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
