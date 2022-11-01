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
	// * order guaranteed never to change

	s := new(scan.R)
	s.B = []byte(`some thing`)
	fmt.Println(s)

}

func ExampleR_Scan() {
	s := new(scan.R)
	s.B = []byte(`foo`)

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
	s := new(scan.R)
	s.B = []byte(`abcdefgh`)
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
	s := new(scan.R)
	s.B = []byte(`foo1234`)

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

func ExampleR_Is() {
	s := new(scan.R)
	s.B = []byte(`foo`)

	s.Scan() // never forget to scan with Is (use Peek otherwise)

	fmt.Println(s.Is("fo"))
	fmt.Println(s.Is("bar"))

	// Output:
	// true
	// false
}

func ExampleR_Peek() {

	s := new(scan.R)
	s.B = []byte(`foo`)

	fmt.Println(s.Peek("fo"))
	s.Scan()
	fmt.Println(s.Peek("fo"))
	fmt.Println(s.Peek("oo"))

	// Output:
	// true
	// false
	// true
}

func ExampleR_Is_not() {
	s := new(scan.R)
	s.B = []byte("\r\n")

	s.Scan() // never forget to scan with Is (use Peek otherwise)

	fmt.Println(s.Is("\r"))
	fmt.Println(s.Is("\r\n"))
	fmt.Println(s.Is("\n"))

	// Output:
	// true
	// true
	// false

}

func ExampleR_Match() {
	s := new(scan.R)
	s.B = []byte(`foo`)

	s.Scan() // never forget to scan (use PeekMatch otherwise)

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

func ExampleR_Pos() {
	s := new(scan.R)
	s.B = []byte("one line\nand another\r\nand yet another")

	s.P = 2
	s.Pos().Print()

	s.P = 0
	s.Scan()
	s.Scan()
	s.Pos().Print()

	s.P = 12
	s.Pos().Print()

	s.P = 27
	s.Pos().Print()

	// Output:
	// U+006E 'n' 1,2-2 (2-2)
	// U+006E 'n' 1,2-2 (2-2)
	// U+0064 'd' 2,3-3 (12-12)
	// U+0079 'y' 3,5-5 (27-27)

}

func ExampleR_Positions() {
	s := new(scan.R)
	s.B = []byte("one line\nand another\r\nand yet another")

	for _, p := range s.Positions(2, 12, 27) {
		p.Print()
	}

	// Output:
	// U+006E 'n' 1,2-2 (2-2)
	// U+0064 'd' 2,3-3 (12-12)
	// U+0079 'y' 3,5-5 (27-27)

}

func ExampleR_Report() {
	defer log.SetFlags(log.Flags())
	defer log.SetOutput(os.Stderr)
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	s := new(scan.R)
	s.B = []byte("one line\nand another\r\nand yet another")

	s.Scan()
	s.Report()

	s.P = 14
	s.Report()

	s.Error("sample error")
	s.Report()

	// Output:
	// U+006F 'o' 1,1-1 (1-1)
	// U+0061 'a' 2,5-5 (14-14)
	// error: sample error at U+0061 'a' 2,5-5 (14-14)

}

func ExampleR_End() {
	s := new(scan.R)
	s.B = []byte("foo")
	s.Print()
	s.Scan()
	s.Print()
	fmt.Println(s.End())
	s.Scan()
	s.Print()
	fmt.Println(s.End())
	s.Scan()
	s.Print()
	fmt.Println(s.End())
	// Output:
	// 0 '\x00' "foo"
	// 1 'f' "oo"
	// false
	// 2 'o' "o"
	// false
	// 3 'o' ""
	// true
}
