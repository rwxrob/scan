package scan_test

import (
	"fmt"
	"regexp"

	"github.com/rwxrob/scan"
)

func ExampleR_init() {

	// * extremely minimal initialization
	// * no need for pointer
	// * order guaranteed never to change

	s := scan.R{B: []byte(`some thing`)}
	fmt.Println(s)

}

func ExampleR_Scan() {
	s := scan.R{B: []byte(`foo`)}

	s.Print() // equivalent of a "zero value"

	fmt.Println(s.Scan())
	s.Print()
	fmt.Println(s.Scan())
	s.Print()
	fmt.Println(s.Scan())
	s.Print()
	fmt.Println(s.Scan()) // does not advance
	s.Print()             // same as before

	// Output:
	// 0 '\x00' "foo"
	// true
	// 1 'f' "oo"
	// true
	// 2 'o' "o"
	// true
	// 3 'o' ""
	// false
	// 3 'o' ""

}

func ExampleR_Scan_loop() {
	s := scan.R{B: []byte(`abcdefgh`)}
	for s.Scan() {
		fmt.Print(string(s.R))
		if s.P != len(s.B) {
			fmt.Print("-")
		}
	}
	// Output:
	// a-b-c-d-e-f-g-h
}

func ExampleR_Scan_jump() {
	s := scan.R{B: []byte(`foo1234`)}

	fmt.Println(s.Scan())
	s.Print()
	s.P += 2
	fmt.Println(s.Scan())
	s.Print()

	// Output:
	// true
	// 1 'f' "oo1234"
	// true
	// 4 '1' "234"

}

func ExampleR_Peek() {
	s := scan.R{B: []byte(`foo`)}

	s.Scan() // never forget to scan (panic otherwise)

	fmt.Println(s.Peek("fo"))
	fmt.Println(s.Peek("bar"))

	// Output:
	// true
	// false
}

func ExampleR_Match() {
	s := scan.R{B: []byte(`foo`)}

	s.Scan() // never forget to scan (panic otherwise)

	f := regexp.MustCompile(`f`)
	F := regexp.MustCompile(`F`)
	o := regexp.MustCompile(`o`)

	fmt.Println(s.Match(f))
	fmt.Println(s.Match(F))
	fmt.Println(s.Match(o))

	// Output:
	// 1
	// -1
	// -1
}
