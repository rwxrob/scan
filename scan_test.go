package scan_test

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rwxrob/scan"
)

func ExampleNew_string() {
	s := scan.New("something here")
	s.Print()
	// Output:
	// U+0073 's' 1,1-1 (1-1)
}

func ExampleNew_runes() {
	s := scan.New([]byte{'f', 'o', 'o'})
	s.Print()
	// Output:
	// U+0066 'f' 1,1-1 (1-1)
}

func ExampleNew_reader() {
	s := scan.New(strings.NewReader("something here"))
	s.Print()
	// Output:
	// U+0073 's' 1,1-1 (1-1)
}

func ExampleR_Init_string() {
	s := new(scan.R)
	s.Init("something here")
	s.Print()
	// Output:
	// U+0073 's' 1,1-1 (1-1)
}

func ExampleR_Init_runes() {
	s := new(scan.R)
	s.Init([]byte{'f', 'o', 'o'})
	s.Print()
	// Output:
	// U+0066 'f' 1,1-1 (1-1)
}

func ExampleR_Init_reader() {
	s := new(scan.R)
	s.Init(strings.NewReader("something here"))
	s.Print()
	// Output:
	// U+0073 's' 1,1-1 (1-1)
}

func ExampleR_Init_error() {
	s := new(scan.R)
	s.Init([]rune{'f', 'o', 'o'})
	s.Print()
	// Output:
	// buffer: unsupported type: []int32 at U+0000 '\x00' 1,1-1 (1-1)
}

func ExampleR_marshaling() {

	// adjust log output for testing
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.Flags())

	s := scan.New("something here")
	s.Print()
	s.Log()
	fmt.Println(s)
	fmt.Println(s.String())

	// Output:
	// U+0073 's' 1,1-1 (1-1)
	// U+0073 's' 1,1-1 (1-1)
	// U+0073 's' 1,1-1 (1-1)
	// U+0073 's' 1,1-1 (1-1)
}

func ExampleR_Scan() {
	s := scan.New("so")
	fmt.Println(s.State == s.State|scan.EOD)
	s.Print()
	s.Scan()
	s.Print()
	s.Scan()
	s.Print()
	fmt.Println(s.State == s.State|scan.EOD)
	// Output:
	// false
	// U+0073 's' 1,1-1 (1-1)
	// U+006F 'o' 1,2-2 (2-2)
	// <EOD>
	// true
}

func ExampleR_Mark() {
	s := scan.New("so")
	s.Print()
	m := s.Mark()
	m.Print()
	s.Scan()
	s.Print()
	m.Print()
	// Output:
	// U+0073 's' 1,1-1 (1-1)
	// U+0073 's' 1,1-1 (1-1)
	// U+006F 'o' 1,2-2 (2-2)
	// U+0073 's' 1,1-1 (1-1)
}

func ExampleR_Snap_one() {
	s := scan.New("something here")
	s.Any(4)
	s.Print()
	s.Snap()
	s.Any(4)
	s.Print()
	s.Back()
	s.Print()
	// Output:
	// U+0074 't' 1,5-5 (5-5)
	// U+0067 'g' 1,9-9 (9-9)
	// U+0074 't' 1,5-5 (5-5)
}

func ExampleR_Snap_nested() {
	s := scan.New("something here")
	s.Any(4)
	s.Print() // t
	s.Snap()  // first time
	s.Any(4)  //
	s.Snap()  // second time
	s.Print() // g
	s.Scan()  //
	s.Print() // ' '
	//s.Snapped.Print()
	s.Back()  // back to second snap
	s.Print() // g
	s.Back()  // back to first snap
	s.Print() // t
	s.Back()  // (does nothing)
	s.Print() // t
	// Output:
	// U+0074 't' 1,5-5 (5-5)
	// U+0067 'g' 1,9-9 (9-9)
	// U+0020 ' ' 1,10-10 (10-10)
	// U+0067 'g' 1,9-9 (9-9)
	// U+0074 't' 1,5-5 (5-5)
	// U+0074 't' 1,5-5 (5-5)
}

func ExamplePeek() {
	s := scan.New("some thing")
	s.Any(6)
	fmt.Println(s.Peek(3))
	// Output:
	// hin
}

func ExampleR_PeekTo() {
	s := scan.New("some thing")
	s.Scan()
	m1 := s.Mark()
	m1.Print()
	s.Any(3)
	fmt.Printf("%q\n", s.PeekTo(m1)) //  behind
	s.Any(4)
	m2 := s.Mark()
	m2.Print()
	s.Jump(m1)
	fmt.Printf("%q\n", s.PeekTo(m2)) //  ahead
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
	// "ome "
	// U+006E 'n' 1,9-9 (9-9)
	// "ome thin"
}

func ExampleR_PeekSlice() {
	s := scan.New("some thing")
	s.Scan()
	m1 := s.Mark()
	m1.Print()
	s.Any(7)
	m2 := s.Mark()
	m2.Print()
	fmt.Println(s.PeekSlice(m1, m2))
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
	// U+006E 'n' 1,9-9 (9-9)
	// ome thin
}
