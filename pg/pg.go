package pg

import (
	"unicode"

	"github.com/rwxrob/scan"
	z "github.com/rwxrob/scan/is"
)

var WS = z.I{SP, TAB, CR, LF}
var EndLine = z.I{LF, CRLF, CR}

func UGraphic(s *scan.R) bool { return unicode.IsGraphic(s.Cur.Rune) }
