# Go Rune Scanner with Cursors

[![Go Version](https://img.shields.io/github/go-mod/go-version/rwxrob/scan)](https://tip.golang.org/doc/go1.18)
[![GoDoc](https://godoc.org/github.com/rwxrob/scan?status.svg)](https://godoc.org/github.com/rwxrob/scan)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)
[![Go Report
Card](https://goreportcard.com/badge/github.com/rwxrob/scan)](https://goreportcard.com/report/github.com/rwxrob/scan)

## Stages of Grammar / Domain-Specific Language Design

1. Work it out with PEGN and scan.X
2. Create Early Tools with scan.X Functions 
3. Refine Grammar Over Time
4. Solidify a PEGN Specification
5. Generate High-Performance Parser

Designing and developing a grammar or domain-specific language is messy
business at first. Flexibility is paramount. Quickly executed
refinements and highly testable/fuzzable code are the priority in the
early stages. PEGN and scan.X work well together to fulfill this need.
PEGN can be used to write the grammar and quickly converted to scan.X
functions directly --- often without even leaving your editor --- and
these functions can immediately be tested and fuzzed to ensure they are
what is wanted. Once the initial grammar scanner/parser is producing the
desired result can be built into alpha versions of the application using
the unoptimized scan.X functions. Over the initial alpha/beta life cycle
of the application refinements can easily be delivered using the scan.X
functions (with PEGN changes as well). Later, after the tools have seen
some heavy use in the beta testing, the scan.X functions (or PEGN) can
be rendered as highly optimized scanner or parser code, with or without
AST support. This approach ensures the fastest design of the grammar
while allowing for maximized optimization later when and if needed.
