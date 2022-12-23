// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package scan provides a rudimentary, performant rune scanner designed
primarily with Parsing Expression Grammars in mind and assumes the
entire buffer to be scanned can be fully loaded into memory and
consisting of nothing but valid unicode code points.  This dramatically
improves scanning speed at the cost of memory required to scan. Consider
the conventional bufio.Scanner as an alternative when memory pressure is
a concern.
*/
package scan

import (
	"io"
	"os"
)

// R (to avoid stuttering) implements a buffered data, non-linear,
// rune-centric, scanner with regular expression support
type R []rune

// Buffer sets Buf and resets Cur to 0. Input may be io.Reader, []byte,
// []rune, or string.  Buffer is typically called immediately after
// instantiating a new scanner (ex: s := new(scan.R).Buffer(in)).
// Returns self-reference.
func (s R) Buffer(in any) R {
	switch v := in.(type) {
	case string:
		s = []rune(v)
	case []byte:
		s = []rune(string(v))
	case []rune:
		s = v
	case io.Reader:
		b, _ := io.ReadAll(v)
		if b != nil {
			s = []rune(string(b))
		}
	}
	return s
}

// Open opens the file at path and loads it by passing to Buffer.
// Returns self-reference.
func (s R) Open(path string) R {
	f, _ := os.Open(path)
	if f == nil {
		return s
	}
	defer f.Close()
	return s.Buffer(f)
}
