# Unicode Codepoint (Rune) Scanner in Go

[![GoDoc](https://godoc.org/github.com/rwxrob/scan?status.svg)](https://godoc.org/github.com/rwxrob/scan)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)
[![Go Report
Card](https://goreportcard.com/badge/github.com/rwxrob/scan)](https://goreportcard.com/report/github.com/rwxrob/scan)

This scanner is meant either to be used as-is or "vendored" and included in other projects to remove any fragile dependencies on what is usually a core software component of project that require such a scanner.

* Fast
* Simple
* Intuitive

***Version v0.2 is a significant breaking change*** since it changes all the methods and the underlying buffer to be `[]rune` instead of `[]byte`.

***Breaking change between v0.10 and v0.11***

