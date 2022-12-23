// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package scan

import (
	"fmt"
	"log"
)

// Position contains the human-friendly information about the position
// within a give text file. In order to do this the full bytes buffer
// must be traversed and lines counted according to the text being
// scanned. Note that all values begin with 1 and not 0. See
// Pointer.Position() and scan.R.Position() for more.
type Position struct {
	Rune    rune // rune at this location
	BufByte int  // byte offset in file
	BufRune int  // rune offset in file
	Line    int  // line offset
	LByte   int  // line column byte offset
	LRune   int  // line column rune offset
}

// String fulfills the fmt.Stringer interface by printing
// the Position in a human-friendly way:
//
//   U+1F47F '👿' 1,3-5 (3-5)
//                | | |  | |
//             line | |  | overall byte offset
//   line rune offset |  overall rune offset
//     line byte offset
//
func (p Position) String() string {
	s := fmt.Sprintf(`%U %q %v,%v-%v (%v-%v)`,
		p.Rune, p.Rune,
		p.Line, p.LRune, p.LByte,
		p.BufRune, p.BufByte,
	)
	return s
}

// Print prints the Position itself in String form. See String.
func (p Position) Print() { fmt.Println(p.String()) }

// Log calls log.Println on Position itself in String form. See String.
func (p Position) Log() { log.Println(p.String()) }
