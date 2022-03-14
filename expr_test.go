package scan_test

import (
	"fmt"

	"github.com/rwxrob/scan"
	z "github.com/rwxrob/scan/is"
	"github.com/rwxrob/scan/tk"
)

func ExampleX_rune() {
	s := scan.New("some thing")
	s.X('s')
	s.Print()
	s.X(tk.ANY)
	s.Print()
	s.X('\t')
	s.Print()
	fmt.Println(s.State == s.State|scan.ERR)
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
	// U+006D 'm' 1,3-3 (3-3)
	// expected '\t' at U+006D 'm' 1,3-3 (3-3)
	// true

}

func ExampleX_string() {
	s := scan.New("some thing")
	s.X("so")
	s.Print()
	s.X("M")
	s.Print()
	// Output:
	// U+006D 'm' 1,3-3 (3-3)
	// expected "M" at U+006D 'm' 1,3-3 (3-3)
}

func ExampleX_in() {
	s := scan.New("some thing")
	s.Scan()
	s.X(z.I{'O', 'o', "ome"})
	s.Print()
	s.X(z.I{'z', 'q'})
	s.Print()
	// Output:
	// U+006D 'm' 1,3-3 (3-3)
	// expected one of ['z' 'q'] at U+006D 'm' 1,3-3 (3-3)
}

func ExampleX_sequence() {
	s := scan.New("some thing")
	s.X(z.X{'s', "om"})
	s.Print()
	s.X(z.X{'e', '\t'})
	s.Print()
	// Output:
	// U+0065 'e' 1,4-4 (4-4)
	// expected '\t' at U+0020 ' ' 1,5-5 (5-5)
}

func ExampleX_endLine() {
	s := scan.New("some\nth\r\ning\rhere")
	s.X("some", z.EndLine, "th", z.EndLine, "ing", z.EndLine, 'h')
	s.Print()
	// Output:
	// U+0065 'e' 1,15-15 (15-15)
}

func ExampleX_optional() {
	s := scan.New("some thing")
	s.X(z.O{'s', 'S'})
	s.Print()
	s.X(z.O{'z', 'x'})
	s.Print()
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
	// U+006F 'o' 1,2-2 (2-2)
}
