package scan_test

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rwxrob/scan"
)

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
	err := s.Init([]rune{'f', 'o', 'o'})
	fmt.Println(err)
	// Output:
	// scanner: unsupported input type: []int32
}

func ExampleR_marshaling() {

	// adjust log output for testing
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.Flags())

	s, _ := scan.New("something here")
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
	s, _ := scan.New("so")
	fmt.Println(s.State == s.State|scan.Done)
	s.Print()
	s.Scan()
	s.Print()
	s.Scan()
	s.Print()
	fmt.Println(s.State == s.State|scan.Done)
	// Output:
	// false
	// U+0073 's' 1,1-1 (1-1)
	// U+006F 'o' 1,2-2 (2-2)
	// <EOD>
	// true
}

func ExampleR_ScanN() {
	s, _ := scan.New("so")
	s.Print()
	s.ScanN(2)
	s.Print()
	// Output:
	// U+0073 's' 1,1-1 (1-1)
	// <EOD>
}

func ExampleR_Mark() {
	s, _ := scan.New("so")
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
