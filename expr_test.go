package scan_test

import (
	"fmt"

	"github.com/rwxrob/scan"
	z "github.com/rwxrob/scan/is"
	"github.com/rwxrob/scan/tk"
)

func ExampleX_rune_Success() {
	s := scan.New("some thing")
	s.X('s')
	s.Print()
	s.X(tk.ANY)
	s.Print()
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
	// U+006D 'm' 1,3-3 (3-3)
}

func ExampleX_rune_Fail() {
	s := scan.New("some thing")
	s.X('S')
	s.Print() // not advanced
	fmt.Println(s.State == s.State|scan.ERR)
	// Output:
	// expected 'S' at U+0073 's' 1,1-1 (1-1)
	// true
}

func ExampleX_in() {
	s := scan.New("some thing")
	s.Scan()
	s.X(z.I{'O', 'o', "ome"})
	s.Print()
	// Output:
	// U+006D 'm' 1,3-3 (3-3)
}
