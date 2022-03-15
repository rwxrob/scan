// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package z ("is") defines the Go scan.X interpreted expression language
(which is passed directly to the scan.X method). The language is
implemented entirely using valid Go types (mostly slices) and
first-class functions of the form func(s *scan.R) bool. Most expressions
will advance the scan to the end of the match but can also be placed
within lookahead expressions (z.Y/z.N) to cancel advancement.
Expressions can be easily combine and combinations can be easily shared
and imported through the use of Go modules and packages. Generating
scan.X Go code from other expression grammars is rather trivial and was
the primary motivation behind the scan.X syntax --- particularly when
a scanner is clearly a better choice over multiple regular expressions,
for example, when creating domain specific languages, grammars, linters,
and language servers. (See PEGN, PEG, ABNF, EBNF, regular expressions,
and others as well as the list of dependent Go modules.)
*/
package z

// P ("parse") is a named sequence of expressions that will be parsed
// and captured as a new Node and added to the scan.R.Nodes field
// effectively turning the scan.R into a parser as well. The first item
// must be an integer (usually a constant) identifying the type of Node.
// If any expression fails to match the scan fails.  Otherwise, a new
// tree.Node[string] is added under the current node and the scan
// proceeds.
type P []any

// X ("expression") is a sequence of expressions used for grouping.  If
// any are not the scan fails. (Equal to (?foo) in regular expressions.)
type X []any

// A ("any") advances exactly v[0] times matching any rune.
type A []int

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

// R ("range","rune") is a advancing expression that
// matches a single Unicode code point (rune, int32) from an inclusive
// consecutive range from first (v[0]) to last (v[1]).
type R []rune

// MM ("minmax") is an advancing expression that matches an inclusive
// minimum (v[0]) and maximum (v[1]) count of the given expression
// (v[2]) in "greedy" fashion (the maximum possible matches are
// advanced).
type MM []any

// M ("min") is an advancing expression that matches an inclusive
// minimum number (v[0]) of the given expression item (v[1]) "greedily".
type M []any

// M0 is shorthand for z.M{0,This}. This is useful to make otherwise
// optional matches "greedy".
type M0 []any

// M1 is shorthand for z.M{1,This}.
type M1 []any

// C is a parameterized advancing expression that matches an exact count
// (v[0]) of the given expression (v[1]).
type C []any
