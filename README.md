# Basic Go Rune Scanner

[![GoDoc](https://godoc.org/github.com/rwxrob/scan?status.svg)](https://godoc.org/github.com/rwxrob/scan)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)
[![Go Report
Card](https://goreportcard.com/badge/github.com/rwxrob/scan)](https://goreportcard.com/report/github.com/rwxrob/scan)

* Fast
* Simple
* Intuitive

## Design Considerations

* **Just a buffer, position, and last rune scanned.**

  There's no need to complicate things for this scanner. This has the
  unexpected benefit of making this scanner something that can be
  marshaled in its current state.

* **Promote easy initialization.**

  Using `s := scan.R{Buf: []byte("foo")}` is a perfect good way to
  initialize this bare-bones scanner and should be promoted in order to
  keep the scope of this scanner very low-level and performant.

* **Trace is an int for later bitwise flags.**

  The benefit of adding levels of tracing are obvious, just not ready to
  add them at the moment.
