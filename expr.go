package scan

import (
	z "github.com/rwxrob/scan/is"
	"github.com/rwxrob/scan/tk"
)

// X is the GOPEGN language interpreter that will process any number of
// valid GOPEGN expressions for scanning, looking ahead, parsing, and
// executing first-class functions. The GOPEGN language can be directly
// generated from any valid PEGN making it ideal for quickly needed
// domain specific languages and other grammars. See the "is" and "tk"
// packages. X will push an Error immediately return if any error is
// encountered. For simplicity, the interpreter uses functional
// recursion in its implementation which might be undesirable in certain
// exceptional cases where grammars produce unusually deep nesting. This
// should not be a concern for most applications.
func (s *R) X(expr ...any) bool {
	s.Snap()

	// same as z.X if more than one
	if len(expr) > 1 {
		for _, r := range expr {
			if !s.X(r) {
				s.Back()
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
		s.Back()
		s.Errorf(`expected %q`, v)
		return false

	case string: // (just a sequence of runes, advances) ----------------
		for _, i := range []rune(v) {
			if s.Cur.Rune != i {
				s.Back()
				s.Errorf(`expected %q`, v)
				return false
			}
			s.Scan()
		}
		return true

	case z.I: // "in" (one of required, advances, in order) -------------
		for _, i := range v {
			s.Snap()
			if s.X(i) {
				return true
			}
			s.Err.Pop()
			s.Back()
		}
		s.Errorf(`expected one of %q`, v)
		return false

	case z.O: // "optional" (if any advances, not required) ------------
		for _, i := range v {
			s.Snap()
			if s.X(i) {
				return true
			}
			s.Err.Pop()
			s.Back()
		}
		return true

	case z.X: // "expression" (each must match in order, advances) -----
		for _, i := range v {
			if !s.X(i) {
				s.Back()
				return false
			}
		}
		return true

	default: // --------------------------------------------------------
		s.Back()
		s.Errorf(`unsupported expression type %T`, v)
		return false

	} // end switch

	return false
}
