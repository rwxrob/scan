package scan

import (
	z "github.com/rwxrob/scan/is"
	"github.com/rwxrob/scan/tk"
)

// X is an expression language interpreter that will process any number
// of valid expressions for scanning, looking ahead, parsing, and
// executing first-class functions. See the "is" and "tk" packages.
// X will push an Error and immediately return if any error is
// encountered. For simplicity, the interpreter uses functional
// recursion in its implementation which might be undesirable in certain
// exceptional cases where grammars produce unusually deep nesting. This
// should not be a concern for most applications.
func (s *R) X(expr ...any) bool {

	// same as z.X if more than one
	if len(expr) > 1 {
		m := s.Mark()
		for _, r := range expr {
			if !s.X(r) {
				s.Jump(m)
				return false
			}
		}
		return true
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
		s.Errorf(`expected %q`, v)
		return false

	case string: // (just a sequence of runes, advances) ----------------
		m := s.Mark()
		for _, i := range []rune(v) {
			if s.Cur.Rune != i {
				s.Jump(m)
				s.Errorf(`expected %q`, v)
				return false
			}
			s.Scan()
		}
		return true

	case z.P: // "parse" (parse tree node) ------------------------------
		cur := s.Nodes.Peek()
		n := s.Tree.Node(v[0].(int), "")
		m := s.Mark()
		r := *s.Tree.Root
		s.Nodes.Push(n)
		defer s.Nodes.Pop()
		if !s.X(v[1:]...) {
			s.Jump(m)
			s.Tree.Root = &r
			return false
		}
		n.V = s.PeekSlice(m, s.Last)
		cur.Append(n)
		return true

	case z.X: // "expression" (each must match in order, advances) ------
		m := s.Mark()
		for _, i := range v {
			if !s.X(i) {
				s.Jump(m)
				return false
			}
		}
		return true

	case z.A: // "any" (advances exactly N of any rune) -----------------
		for i := 0; i < v[0]; i++ {
			s.Scan()
		}
		return true

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
				s.Errorf(`unexpected %q`, i)
				s.Jump(m)
				return false
			}
			s.Err.Pop()
		}
		s.Jump(m)
		return true

	case z.I: // "in" (one of required, advances, ordered) --------------
		for _, i := range v {
			m := s.Mark()
			if s.X(i) {
				return true
			}
			s.Err.Pop()
			s.Jump(m)
		}
		s.Errorf(`expected one of %q`, v)
		return false

	case z.O: // "optional" (if any advances, not required) -------------
		for _, i := range v {
			m := s.Mark()
			if s.X(i) {
				return true
			}
			s.Err.Pop()
			s.Jump(m)
		}
		return true

	case z.T: // "to inclusive" (advances to match and includes) -------
		m := s.Mark()
		for s.Cur.Rune != tk.EOD {
			m := s.Mark()
			if s.X(z.X(v)) {
				s.Jump(m)
				return true
			}
			s.Err.Pop()
			s.Scan()
		}
		s.Jump(m)
		s.Errorf("%q not found anywhere in remaining buffer starting", v)
		return false

	case z.Ti: // "to" (advances to match and excludes) ------------------
		m := s.Mark()
		for s.Cur.Rune != tk.EOD {
			if s.X(z.X(v)) {
				return true
			}
			s.Err.Pop()
			s.Scan()
		}
		s.Jump(m)
		s.Errorf("%q not found anywhere in remaining buffer starting", v)
		return false

	case z.R: // "range" (inclusive range between rune int values) ------
		if v[0] <= s.Cur.Rune && s.Cur.Rune <= v[1] {
			s.Scan()
			return true
		}
		s.Errorf(`expected %q-%q`, v[0], v[1])
		return false

	case z.MM: // "min max" (minimum and maximum count of, advances) ----
		count := 0
		for s.Cur.Rune != tk.EOD {
			if !s.X(v[2]) {
				s.Err.Pop()
				break
			}
			count++
		}
		if v[0].(int) <= count && count <= v[1].(int) {
			return true
		}
		s.Errorf(`expected %v-%v of %q`, v[0], v[1], v[2])
		return false

	case z.M: // "min" (minimum and maximum count of, advances) ---------
		count := 0
		for s.Cur.Rune != tk.EOD {
			if !s.X(v[1]) {
				s.Err.Pop()
				break
			}
			count++
		}
		if v[0].(int) <= count {
			return true
		}
		s.Errorf(`expected at least %v of %q`, v[0], v[1])
		return false

	case z.M0: // "min zero" (shorthand for z.M{0,This}) -----------------
		return s.X(z.M{0, v[0]})

	case z.M1: // "min one" (shorthand for z.M{1,This}) -----------------
		return s.X(z.M{1, v[0]})

	case z.C: // "count" (match exactly N of, advances) -----------------
		m := s.Mark()
		for i := 0; i < v[0].(int); i++ {
			if !s.X(v[1]) {
				s.Jump(m)
				s.Errorf(`expected %v of %q`, v[0].(int), v[1])
				return false
			}
		}
		return true

	case func(s *R) bool:
		return v(s)

	default: // ---------------------------------------------------------
		s.Errorf(`unsupported expression type %T`, v)
		return false

	} // end switch

	return false
}
