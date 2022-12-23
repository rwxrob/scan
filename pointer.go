// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package scan

import (
	"fmt"
)

var PointerView = 10 // default length of preview window, 0 disables

// Pointer contains a pointer to a bytes slice buffer and information
// pointing to a specific rune within that buffer. This minimal
// information is required for for any rune scanner since the width
// in bytes of a rune varies between one and four. The order of
// fields is guaranteed to never change (ex: Pointer{buf,'x',4,5}). For
// more contextual positional information see Position.
type Pointer struct {
	Buf *[]byte // pointer to actual bytes buffer
	R   rune    // last rune scanned
	B   int     // beginning of last rune scanned
	E   int     // effective end of last rune scanned (beginning of next)
}

// String implements fmt.Stringer with the last rune scanned (R),
// and the beginning and ending byte positions joined with
// a dash along with PointerView length of upcoming bytes for context.
// For better positional information in output use Position instead.
func (p Pointer) String() string {
	if PointerView > 0 {
		end := p.E + PointerView
		if end > len((*p.Buf)) {
			end = len((*p.Buf))
		}
		return fmt.Sprintf("%q %v-%v  %q", p.R, p.B, p.E, (*p.Buf)[p.E:end])
	}
	return fmt.Sprintf("%q %v-%v", p.R, p.B, p.E)
}
