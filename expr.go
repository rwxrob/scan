package scan

import (
	"log"

	z "github.com/rwxrob/scan/is"
	"github.com/rwxrob/scan/tk"
	"github.com/rwxrob/to"
)

// X is an moderately performant, easy to use, expression language
// interpreter and rooted node tree parser that will process any number
// of valid expressions for scanning, looking ahead, parsing, and
// executing first-class functions. The resulting parse tree, marshaled
// as JSON, is far more useful than the output of regular expressions.
// It's purpose is to facilitate rapid grammar creation, grammar
// transcription, and easy code generation for more optimized parsers
// when needed. Emphasis has been placed on speed of development over
// raw code execution. As such the interpreter and parser employ
// functional recursion and cached node-tree states which might be
// undesirable in certain exceptional cases where grammars produce
// unusually deep nesting. This should not be a concern for most
// applications where scan.X performance is more than sufficient ---
// particularly for Go applications that replace casual shell scripts.
// See the "is" and "tk" packages for the language specification. X will
// push an Error and immediately return if any error is encountered.
func (s *R) X(expr ...any) bool {

	// same as z.X if more than one
	if len(expr) > 1 {
		m := s.Mark()
		save := s.Tree.Root.Copy()
		for _, r := range expr {
			if !s.X(r) {
				s.Tree.Root = save
				s.Jump(m)
				return false
			}
		}
		return true
	}

	if s.tracex {
		s.Log()
		log.Print(to.Human(expr[0]))
		log.Print("---------------")
	}

	switch v := expr[0].(type) {

	case rune: // -------------------------------------------------------
		if s.Cur.Rune == v || v == tk.ANY {
			s.Scan()
			return true
		}
		if v == tk.NL {
			s.Cur.NewLine()
			return true
		}
		s.Errorf(v, `scan.x: expected %q`, v)
		return false

	case string: // (just a sequence of runes, advances) ----------------
		m := s.Mark()
		for _, i := range []rune(v) {
			if s.Cur.Rune != i {
				s.Jump(m)
				s.Errorf(v, `scan.x: expected %q`, v)
				return false
			}
			s.Scan()
		}
		return true

	case z.X: // "expression" (each must match in order, advances) ------
		m := s.Mark()
		save := s.Tree.Root.Copy()
		for _, i := range v {
			if !s.X(i) {
				s.Tree.Root = save
				s.Jump(m)
				return false
			}
		}
		return true

	case z.A: // "any" (advances exactly N of any rune) -----------------
		return s.Any(v.N)

	case z.Y: // "yes" (positive look-ahead, no advance, ordered) -------
		m := s.Mark()
		if s.X(z.X(v)) {
			s.Jump(m)
			return true
		}
		s.Jump(m)
		return false

	case z.N: // "no" "neg" (negative look-ahead, no advance, ordered) --
		m := s.Mark()
		for _, i := range v {
			if s.X(i) {
				s.Errorf(v, `scan.x: unexpected %v`, to.Human(i))
				s.Jump(m)
				return false
			}
			s.ClearLastError()
		}
		s.Jump(m)
		return true

	case z.I: // "in" (one of required, advances, ordered) --------------
		m := s.Mark()
		save := s.Tree.Root.Copy()
		for _, i := range v {
			if s.X(i) {
				return true
			}
			s.ClearLastError()
			s.Jump(m)
		}
		s.Tree.Root = save
		s.Errorf(v, `scan.x: expected %v`, to.Human(v))
		return false

	case z.O: // "optional" (if any advances, not required) -------------
		m := s.Mark()
		for _, i := range v {
			if s.X(i) {
				return true
			}
			s.Err.Pop()
			s.Jump(m)
		}
		return true

	case z.T: // "to inclusive" (advances to match and includes) -------
		m := s.Mark()
		for {
			m := s.Mark()
			if s.X(v.This) {
				s.Jump(m)
				return true
			}
			s.ClearLastError()
			if !s.Scan() {
				break
			}
		}
		s.Jump(m)
		s.Errorf(v, "scan.x: %v not found", v)
		return false

	case z.Ti: // "to" (advances to match and excludes) ------------------
		m := s.Mark()
		for {
			if s.X(v.This) {
				return true
			}
			s.ClearLastError()
			if !s.Scan() {
				break
			}
		}
		s.Jump(m)
		s.Errorf(v, "scan.x: %v not found", v)
		return false

	case z.R: // "range" (inclusive range between rune int values) ------
		m := s.Mark()
		if v.First <= s.Cur.Rune && s.Cur.Rune <= v.Last {
			s.Scan()
			return true
		}
		s.Jump(m)
		s.Errorf(v, `scan.x: expected %v`, v)
		return false

	case z.MM: // "min max" (minimum and maximum count of, advances) ----
		m := s.Mark()
		count := 0
		for s.Cur.Rune != tk.EOD {
			if !s.X(v.This) {
				s.ClearLastError()
				break
			}
			count++
		}
		if v.Min <= count && count <= v.Max {
			return true
		}
		s.Jump(m)
		s.Errorf(v, `scan.x: expected %v`, to.Human(v))
		return false

	case z.M: // "min" (minimum and maximum count of, advances) ---------
		m := s.Mark()
		save := s.Tree.Root.Copy()
		count := 0
		for s.Cur.Rune != tk.EOD {
			if !s.X(v.This) {
				s.ClearLastError()
				break
			}
			count++
		}
		if v.Min <= count {
			return true
		}
		s.Tree.Root = save
		s.Jump(m)
		s.Errorf(v, `scan.x: expected %v`, to.Human(v))
		return false

	case z.M0: // "min zero" (shorthand for z.M{0,This}) -----------------
		return s.X(z.M{0, v.This})

	case z.M1: // "min one" (shorthand for z.M{1,This}) -----------------
		return s.X(z.M{1, v.This})

	case z.C: // "count" (match exactly N of, advances) -----------------
		m := s.Mark()
		for i := 0; i < v.N; i++ {
			if !s.X(v.This) {
				s.Jump(m)
				s.Errorf(v, `scan.x: expected %v`, to.Human(v))
				return false
			}
		}
		return true

	case func(s *R) bool: // first-class function hook (does whatever) --
		return v(s)

	case z.P: // "parse" (parse tree node) ------------------------------
		cur := s.Nodes.Peek()
		n := s.Tree.Node(v.T, "")
		m := s.Mark()
		s.Nodes.Push(n)
		defer s.Nodes.Pop()
		if !s.X(v.This) {
			s.Jump(m)
			return false
		}
		n.V = s.PeekSlice(m, s.Last)
		cur.Append(n)
		return true

	default: // ---------------------------------------------------------
		s.Errorf(v, `scan.x: unsupported %T`, v)
		return false

	} // end switch

	return false
}
