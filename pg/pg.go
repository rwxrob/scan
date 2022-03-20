package pg

import (
	"unicode"

	"github.com/rwxrob/scan"
)

func UGraphic(s *scan.R) bool { return unicode.IsGraphic(s.Cur.Rune) }
