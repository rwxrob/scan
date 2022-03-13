// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package z (often imported as "is") defines the GOPEGN language
implemented entirely using Go types (mostly slices and structs). GOPEGN
can be 100% transpiled to and from the Parsing Expression Grammer
Notation (PEGN).

Slices represent sets of possibilities.

Structs provide parameters for more complex expressions and are are
guaranteed never to change allowing them to be dependably used in
assignment without struct field names using Go's inline composable
syntax. Some editors may need configuring to allow this since in general
practice this can create subtle (but substantial) foot-guns for
maintainers.

"Advancing" expressions will advance the scan to the end of the
expression match.

"Look-ahead" expressions simply check for a match but do not advance the
scan. Developers should take careful note of the difference in the
documentation.

Composites are compound expressions composed of others. They represent
the tokens and classes from PEGN and other grammars and are designed to
simplify grammar development at a higher level.

First-class functions are not strictly an expression type but are fully
supported by the GOPEGN grammar and equate to StateDef PEGN definitions
which are identified by all capitals but begin with underscore (_).
GOPEGN generators may choose to create stub functions for such when
transpiling from PEGN and developers may choose to use the same function
naming convention to distinguish such functions even when no code
generation is involved.
*/
package z

// ------------------------------- core -------------------------------

// P ("parse") is a named sequence of expressions that will be parsed
// and captured as a new Node and added to the scan.R.Nodes field
// effectively turning the scan.R into a parser as well. The first item
// must be an integer (usually a constant) identifying the type of Node.
// If any expression fails to match the scan fails.  Otherwise, a new
// tree.Node[string] is added under the current node and the scan
// proceeds. Nodes must only contain other nodes or a string value,
// never both. If the first item in the sequence after the type is not
// also a node (z.P) then the node is marked as "edge" (or "leaf") and
// any nodes detected further in the sequence will cause the scan to
// fail with a syntax error.
type P []any

// X ("expression") is a sequence of expressions used for grouping.  If
// any are not the scan fails. (Equal to (?foo) in regular expressions.)
type X []any

// ------------------------------- sets -------------------------------

// Y ("yes") is a set of positive lookahead expressions. If any are
// seen at the current cursor position the scan will proceed without
// consuming them (unlike z.O and z.I). If none are found the scan
// fails. (Equal to ampersand (&) in PEGN.)
type Y []any

// N ("not") is a set of negative lookahead expressions. If any are seen
// at the current cursor position the scan will fail and the scan is
// never advanced. This is useful when everything from one expression is
// wanted except for a few negative exceptions. (Equal to exclamation
// point (!) in PEGN.)
type N []any

// I ("in","include") is a set of advancing expressions. If any
// expression in the slice is found the scan advances to the end of that
// expression and continues. If none of the expressions is found the
// scan fails.  Evaluation of expressions is always left to right
// allowing parser developers to prioritize common expressions at the
// beginning of the slice.
type I []any

// O ("optional") is a set of optional advancing expressions. If any
// expression is found the scan is advanced (unlike is.Y, which does
// not advance).
type O []any

// T ("to") is a set of advancing expressions that mark an exclusive
// boundary at which the scan should stop without including the
// boundary.
type T []any

// Ti ("to inclusive") is an inclusive version of z.T which includes the
// boundary.
type Ti []any

// --------------------------- parameterized --------------------------

// MM ("minmax") is a parameterized advancing expression that matches an
// inclusive minimum and maximum count of the given expression (This).
type MM struct {
	Min  int
	Max  int
	This any
}

// M ("min") is a parameterized advancing expression that matches an
// inclusive minimum number of the given expression item (This). Use
// within is.It to disable advancement.
type M struct {
	Min  int
	This any
}

// M1 is shorthand for z.M{1,This}.
type M1 struct{ This any }

// C is a parameterized advancing expression that matches an exact count
// of the given expression (This). Use within is.It to disable
// advancement.
type C struct {
	N    int
	This any
}

// C2 is shorthand for z.C{2,This}.
type C2 struct{ This any }

// C3 is shorthand for z.C{3,This}.
type C3 struct{ This any }

// C4 is shorthand for z.C{4,This}.
type C4 struct{ This any }

// C5 is shorthand for z.C{5,This}.
type C5 struct{ This any }

// C6 is shorthand for z.C{6,This}.
type C6 struct{ This any }

// C7 is shorthand for z.C{7,This}.
type C7 struct{ This any }

// C8 is shorthand for z.C{8,This}.
type C8 struct{ This any }

// C9 is shorthand for z.C{9,This}.
type C9 struct{ This any }

// A ("any") is short for z.C{N,tk.ANY}.
type A struct {
	N int
}

// R ("range","rune") is a parameterized advancing expression that
// matches a single Unicode code point (rune, int32) from an inclusive
// consecutive set from First to Last (First,Last).
type R struct {
	First rune
	Last  rune
}
