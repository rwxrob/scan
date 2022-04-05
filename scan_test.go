package scan_test

import (
	"fmt"
	"log"
	"os"
	"regexp"

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

func ExampleR_Scan_jump() {
	s := scan.R{Buf: []byte(`foo1234`)}

	fmt.Println(s.Scan())
	s.Print()
	s.Pos += 2
	fmt.Println(s.Scan())
	s.Print()

	// Output:
	// true
	// 1 U+0066 'f'
	// true
	// 4 U+0031 '1'

}

func ExampleR_Is() {
	s := scan.R{Buf: []byte(`foo`)}

	fmt.Println(s.Is("fo"))
	fmt.Println(s.Is("bar"))

	// Output:
	// true
	// false
}

func ExampleR_Match() {
	s := scan.R{Buf: []byte(`foo`)}
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

func ExampleR_View() {
	defer log.SetOutput(os.Stderr)
	log.SetOutput(os.Stdout)
	defer log.SetFlags(log.Flags())
	log.SetFlags(0)

	s := scan.R{Buf: []byte(`foo`)}
	s.View()
	scan.ViewLen = 2
	s.View()

	// Output:
	// "foo"
	// "fo"...
}
