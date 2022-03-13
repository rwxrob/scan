package scan

import (
	z "github.com/rwxrob/scan/is"
	"github.com/rwxrob/scan/tk"
	"github.com/rwxrob/structs/qstack"
)

// X is the GOPEGN language interpreter that will process any number of
// valid GOPEGN expressions for scanning, looking ahead, parsing, and
// executing first-class functions. The GOPEGN language can be directly
// generated from any valid PEGN making it ideal for quickly needed
// domain specific languages and other grammars. See the "is" and "tk"
// packages. X will push an Error immediately return if any error is
// encountered.
func (s *R) X(expr ...any) bool {

	q := qstack.New[any]()
	q.Push(expr...)

	s.Snap()

	for q.Len > 0 {

		switch v := q.Pop().(type) {

		case rune: // ------------------------------------------------------
			if s.Cur.Rune == v || v == tk.ANY {
				s.Scan()
				return true
			}
			s.Back()
			s.Errorf(`expected %q`, v)
			return false

		case z.I: // -----------------------------------------------------
			for _, i := range v {
				s.Snap()
				if s.X(i) {
					return true
				}
				s.Err.Pop()
				s.Back()
			}
			return false

		default:
			s.Back()
			s.Errorf(`unsupported expression type %T`, v)
			return false

		} // end switch

	} // end for

	return false
}
