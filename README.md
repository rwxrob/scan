# Unicode Codepoint (Rune) Scanner in Go

[![GoDoc](https://godoc.org/github.com/rwxrob/scan?status.svg)](https://godoc.org/github.com/rwxrob/scan)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)
[![Go Report
Card](https://goreportcard.com/badge/github.com/rwxrob/scan)](https://goreportcard.com/report/github.com/rwxrob/scan)

This "scanner" is a simple demonstration of how casting a `[]byte` slice to the `[]rune` slice is all that is effectively needed for creating recursive descent parsers.

***Version v0.2 is a significant breaking change*** since it completely drops all the unnecessary complexity of what should have always just been a `[]rune` slice. 

***Breaking change between v0.10 and v0.11***

