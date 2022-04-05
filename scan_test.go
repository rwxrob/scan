package scan_test

import (
	"fmt"

	"github.com/rwxrob/scan"
)

func ExampleR_init() {

	// * extremely minimal initialization
	// * no need for pointer
	// * order guaranteed never to change

	s := scan.R{Buf: []byte(`some thing`)}
	fmt.Println(s)

}

func ExampleR_Scan() {
	s := scan.R{Buf: []byte(`foo`)}

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
	// 0 U+0000 '\x00'
	// true
	// 1 U+0066 'f'
	// true
	// 2 U+006F 'o'
	// true
	// 3 U+006F 'o'
	// false
	// 3 U+006F 'o'
}

func ExampleR_ScanN() {
	s := scan.R{Buf: []byte(`foo`)}

	fmt.Println(s.ScanN(3))
	s.Print()
	fmt.Println(s.ScanN(1))
	s.Print()

	// Output:
	// true
	// 3 U+006F 'o'
	// false
	// 3 U+006F 'o'
}

func ExampleR_Is() {
	s := scan.R{Buf: []byte(`foo`)}

	fmt.Println(s.Is("fo"))
	fmt.Println(s.Is("bar"))

	// Output:
	// true
	// false
}
