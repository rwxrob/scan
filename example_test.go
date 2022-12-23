package scan_test

import (
	"fmt"

	"github.com/rwxrob/scan"
)

func ExampleR() {

	s := new(scan.R).Buffer(`😊 lol`)
	for i, r := range s {
		fmt.Printf("%v %q %v\n", i, r, r)
	}

	// Output:
	// 0 '😊' 128522
	// 1 ' ' 32
	// 2 'l' 108
	// 3 'o' 111
	// 4 'l' 108

}

func ExampleR_recursive_Descent() {

	/*
		Given the following PEGN grammar

		Doc   <-- ws* Title Para ws*
		Title <-- '#' SP rune{1,70} EOL
		EOL    <- CR? NL

	*/

	md := `
# 😊 Title here

Something here that is a **term** to be parsed
and another line.
`

	s := new(scan.R).Buffer(md)

	for i := 0; i < len(s); i++ {

	}

	// Output:
	// some

}
