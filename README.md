# Basic Go Rune Scanner

**💥 DEPRECATED. Moved to
[rwxrob/pegn](https://github.com/rwxrob/pegn).**

[![GoDoc](https://godoc.org/github.com/rwxrob/scan?status.svg)](https://godoc.org/github.com/rwxrob/scan)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)
[![Go Report
Card](https://goreportcard.com/badge/github.com/rwxrob/scan)](https://goreportcard.com/report/github.com/rwxrob/scan)

* Fast
* Simple
* Intuitive

## Design Considerations

**Moved to github.com/rwxrob/pegn/scan**

This scanner is heavily influenced by PEG and expectations about use of
memory and such are better set by moving it into a PEG-related package.

**Breaking change between v0.10 and v0.11**

After adding `pegn.Scanner` interface implementation went ahead and
changed `scan.LP` to `scan.PP` (for previous instead of last).

**Fulfills pegn.Scanner interface**

The `pegn.Scanner` interface has become authoritative for this (and
most) PEGN parsers I develop since it can allow the required
abstractions when necessary without getting in the way of performant
parsing when needed.

**Just a buffer, position, and last rune scanned**

There's no need to complicate things for this scanner. This has the
unexpected benefit of making this scanner something that can be
marshaled in its current state.

* **Promote easy initialization.**

  Using `s := scan.R{B: []byte("foo")}` is a perfect good way to
  initialize this bare-bones scanner and should be promoted in order to
  keep the scope of this scanner very low-level and performant.

* **Trace is an int for later bitwise flags.**

  The benefit of adding levels of tracing are obvious, just not ready to
  add them at the moment.
