package scan_test

import (
	"fmt"
	"log"
	"os"

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
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
	// U+006D 'm' 1,3-3 (3-3)
	// expected '\t' at U+006D 'm' 1,3-3 (3-3)
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

func ExampleX_end_of_Line() {
	eol := z.I{'\n', "\r\n", '\r'}
	s := scan.New("some\nth\r\ning\rhere")
	s.X("some", eol, "th", eol, "ing", eol, 'h')
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

func ExampleX_any() {
	s := scan.New("some thing")
	s.X(z.A{3})
	s.Print()
	s.X(z.A{30})
	s.Print()
	// Output:
	// U+0065 'e' 1,4-4 (4-4)
	// <EOD>
}

func ExampleX_count() {
	s := scan.New("sssome thing")
	s.X(z.C{3, 's'})
	s.Print()
	s.X(z.C{2, 'o'})
	s.Print()
	// Output:
	// U+006F 'o' 1,4-4 (4-4)
	// expected 2 of 'o' at U+006F 'o' 1,4-4 (4-4)
}

func ExampleX_positive_Lookahead() {
	s := scan.New("some thing")
	s.X(z.Y{'s', "om"})
	s.Print()
	s.X(z.Y{'s', 'O'})
	s.Print()
	s.Err.Pop()
	s.X(z.Y{'O'})
	s.Print()
	// Output:
	// U+0073 's' 1,1-1 (1-1)
	// expected 'O' at U+006F 'o' 1,2-2 (2-2)
	// expected 'O' at U+0073 's' 1,1-1 (1-1)
}

func ExampleX_negative_Lookahead() {
	s := scan.New("some thing")
	s.X(z.N{'z'})
	s.Print()
	s.X(z.N{'s'})
	s.Print()
	// Output:
	// U+0073 's' 1,1-1 (1-1)
	// unexpected 's' at U+006F 'o' 1,2-2 (2-2)
}

func ExampleX_to() {
	s := scan.New("some thing")
	s.X(z.T{' '})
	s.Print()
	// Output:
	// U+0020 ' ' 1,5-5 (5-5)
}

func ExampleX_to_Inclusive() {
	s := scan.New("some thing")
	s.X(z.Ti{' '})
	s.Print()
	s.X(z.Ti{'z'})
	s.Print()
	// Output:
	// U+0074 't' 1,6-6 (6-6)
	// ['z'] not found anywhere in remaining buffer starting at U+0074 't' 1,6-6 (6-6)
}

func ExampleX_range() {
	s := scan.New("some thing")
	s.X(z.R{'a', 'z'})
	s.Print()
	s.X(z.R{'A', 'Z'})
	s.Print()
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
	// expected 'A'-'Z' at U+006F 'o' 1,2-2 (2-2)
}

func ExampleX_min_Max() {
	s := scan.New("  sssome thing")
	s.X(z.MM{1, 3, ' '})
	s.Print()
	s.X(z.MM{4, 6, 's'})
	s.Print()
	// Output:
	// U+0073 's' 1,3-3 (3-3)
	// expected 4-6 of 's' at U+006F 'o' 1,6-6 (6-6)
}

func ExampleX_min() {
	s := scan.New("  sssome thing")
	s.X(z.M{1, ' '})
	s.Print()
	s.X(z.M{4, 's'})
	s.Print()
	// Output:
	// U+0073 's' 1,3-3 (3-3)
	// expected at least 4 of 's' at U+006F 'o' 1,6-6 (6-6)
}

func ExampleX_min_One() {
	s := scan.New("  sssome thing")
	s.X(z.M1{' '})
	s.Print()
	s.X(z.M1{'a'})
	s.Print()
	// Output:
	// U+0073 's' 1,3-3 (3-3)
	// expected at least 1 of 'a' at U+0073 's' 1,3-3 (3-3)
}

func ExampleX_first_Class_Functions() {

	// adjust log output for testing
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.Flags())

	logit := func(s *scan.R) bool { s.Log(); return true }
	scanSome := func(s *scan.R) bool { return s.X("some") }
	scanTh := func(s *scan.R) bool { return s.X("th") }
	ws := func(s *scan.R) bool { return s.X(z.I{' ', '\t', '\r', '\n'}) }
	// ws := z.X(' ', '\t', '\r', '\n')

	s := scan.New("some\nthing")
	s.X(scanSome, logit, ws, tk.NL, scanTh)
	s.Print()
	// Output:
	// U+000A '\n' 1,5-5 (5-5)
	// U+0069 'i' 2,3-3 (8-8)
}

func ExampleX_parse_Single() {
	s := scan.New("some thing")
	s.X(z.P{1, 's', 'o', "me"})
	s.Print()
	fmt.Println(s.Tree.Root.Count)
	s.Tree.Root.Print()
	// Output:
	// U+0020 ' ' 1,5-5 (5-5)
	// 1
	// {"T":1,"N":[{"T":1,"V":"some"}]}
}

func ExampleX_parse_Two() {
	s := scan.New("some thing")
	s.X(z.P{2, "so"}, z.P{2, "me"})
	s.Print()
	fmt.Println(s.Tree.Root.Count)
	s.Tree.Root.Print()
	// Output:
	// U+0020 ' ' 1,5-5 (5-5)
	// 2
	// {"T":1,"N":[{"T":2,"V":"so"},{"T":2,"V":"me"}]}

}

func ExampleX_parse_Nested_Simple() {
	s := scan.New("some thing")
	s.X(z.P{2, z.P{3, "some"}})
	s.Print()
	fmt.Println(s.Tree.Root.Count)
	s.Tree.Root.Print()
	// Output:
	// U+0020 ' ' 1,5-5 (5-5)
	// 1
	// {"T":1,"N":[{"T":2,"V":"some","N":[{"T":3,"V":"some"}]}]}

}

func ExampleX_parse_Nested_with_Other() {
	s := scan.New("some thing")
	s.X(z.P{2, "some", ' ', z.P{3, "th"}})
	s.Print()
	fmt.Println(s.Tree.Root.Count)
	s.Tree.Root.Print()
	// Output:
	// U+0069 'i' 1,8-8 (8-8)
	// 1
	// {"T":1,"N":[{"T":2,"V":"some th","N":[{"T":3,"V":"th"}]}]}
}

func ExampleX_parse_Nested_with_Two_Other() {
	s := scan.New("some thing")
	s.X(z.P{2, z.P{3, "some"}, ' ', z.P{3, "th"}})
	s.Print()
	fmt.Println(s.Tree.Root.Count)
	s.Tree.Root.Print()
	// Output:
	// U+0069 'i' 1,8-8 (8-8)
	// 1
	// {"T":1,"N":[{"T":2,"V":"some th","N":[{"T":3,"V":"some"},{"T":3,"V":"th"}]}]}
}

func ExampleX_parse_Nested_Expression() {

	const WORD = 2
	const CHAR = 3

	ch := z.P{CHAR, z.R{'a', 'z'}}
	word := z.P{WORD, z.M1{ch}}
	ws := z.I{' ', '\t', '\r', '\n'}

	s := scan.New("meE")
	s.X(z.X{word, ws}, word)
	s.Print()
	s.Tree.Root.Print()

	// Output:
	// expected one of [' ' '\t' '\r' '\n'] at U+0045 'E' 1,3-3 (3-3)
	// {"T":1}
}

/*
func ExampleX_parse_Nested_Expression_EOD() {

	const WORD = 2
	const CHAR = 3

	ch := z.P{CHAR, z.R{'a', 'z'}}
	word := z.P{WORD, z.M1{ch}}
	ws := z.I{' ', '\t', '\r', '\n'}

	s := scan.New("me")
	s.X(z.X{word, ws}, word)
	s.Print()
	s.Tree.Root.Print()

	// Output:
	// expected one of [' ' '\t' '\r' '\n'] at U+0045 'E' 1,3-3 (3-3)
	// {"T":1}
}
*/

func ExampleX_parse_Nested_Complex() {

	const WORD = 2
	const CHAR = 3

	ch := z.P{CHAR, z.R{'a', 'z'}}
	word := z.P{WORD, z.M1{ch}}
	ws := z.I{' ', '\t', '\r', '\n'}

	s := scan.New("go me again")
	s.X(z.X{word, ws}, word)
	s.Print()
	nodes := s.Tree.Root.Nodes()
	fmt.Println(s.Tree.Root.Count, nodes[0].Count, nodes[1].Count)
	s.Tree.Root.Print()

	// Output:
	// U+0020 ' ' 1,6-6 (6-6)
	// 2 2 2
	// {"T":1,"N":[{"T":2,"V":"go","N":[{"T":3,"V":"g"},{"T":3,"V":"o"}]},{"T":2,"V":"me","N":[{"T":3,"V":"m"},{"T":3,"V":"e"}]}]}
}
