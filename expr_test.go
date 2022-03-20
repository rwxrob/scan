package scan_test

import (
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
	s.X('\t')
	s.Print()
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
	// scan.x: expected '\t' at U+006F 'o' 1,2-2 (2-2)
}

func ExampleX_rune_Any() {
	s := scan.New("some thing")
	s.X(tk.ANY)
	s.Print()
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
}

func ExampleX_rune_NL() {
	s := scan.New("some\nthing")
	s.X("some", '\n', tk.NL, "th")
	s.Print()
	// Output:
	// U+0069 'i' 2,3-3 (8-8)

}

func ExampleX_string() {
	s := scan.New("some thing")
	s.X("so")
	s.Print()
	s.X("M")
	s.Print()
	// Output:
	// U+006D 'm' 1,3-3 (3-3)
	// scan.x: expected "M" at U+006D 'm' 1,3-3 (3-3)
}

func ExampleX_sequence() {
	s := scan.New("some thing")
	s.X(z.X{'s', "om"})
	s.Print()
	s.X(z.X{'e', '\t'})
	s.Print()
	// Output:
	// U+0065 'e' 1,4-4 (4-4)
	// scan.x: expected '\t' at U+0020 ' ' 1,5-5 (5-5)
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
	// scan.x: expected 'O' at U+006F 'o' 1,2-2 (2-2)
	// scan.x: expected 'O' at U+0073 's' 1,1-1 (1-1)
}

func ExampleX_negative_Lookahead() {
	s := scan.New("some thing")
	s.X(z.N{'z'})
	s.Print()
	s.X(z.N{'s'})
	s.Print()
	// Output:
	// U+0073 's' 1,1-1 (1-1)
	// scan.x: unexpected 's' at U+006F 'o' 1,2-2 (2-2)
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
	// scan.x: expected z.I{'z','q'} at U+006D 'm' 1,3-3 (3-3)
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

func ExampleX_to() {
	s := scan.New("some thing")
	s.X(z.T{' '})
	s.Print()
	s.X(z.T{'z'})
	s.Print()
	// Output:
	// U+0020 ' ' 1,5-5 (5-5)
	// scan.x: z.T{'z'} not found at U+0020 ' ' 1,5-5 (5-5)
}

func ExampleX_to_Inclusive() {
	s := scan.New("some thing")
	s.X(z.Ti{' '})
	s.Print()
	s.X(z.Ti{'z'})
	s.Print()
	// Output:
	// U+0074 't' 1,6-6 (6-6)
	// scan.x: z.Ti{'z'} not found at U+0074 't' 1,6-6 (6-6)
}

func ExampleX_range() {
	s := scan.New("some thing")
	s.X(z.R{'a', 'z'})
	s.Print()
	s.X(z.R{'A', 'Z'})
	s.Print()
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
	// scan.x: expected z.R{'A','Z'} at U+006F 'o' 1,2-2 (2-2)
}

func ExampleX_min_Max() {
	s := scan.New("  sssome thing")
	s.X(z.MM{1, 3, ' '})
	s.Print()
	s.X(z.MM{4, 6, 's'})
	s.Print()
	// Output:
	// U+0073 's' 1,3-3 (3-3)
	// scan.x: expected z.MM{4,6,'s'} at U+0073 's' 1,3-3 (3-3)
}

func ExampleX_min() {
	s := scan.New("  sssome thing")
	s.X(z.M{1, ' '})
	s.Print()
	s.X(z.M{4, 's'})
	s.Print()
	// Output:
	// U+0073 's' 1,3-3 (3-3)
	// scan.x: expected z.M{4,'s'} at U+0073 's' 1,3-3 (3-3)
}

func ExampleX_min_One() {
	s := scan.New("  sssome thing")
	s.X(z.M1{' '})
	s.Print()
	s.X(z.M1{'a'})
	s.Print()
	// Output:
	// U+0073 's' 1,3-3 (3-3)
	// scan.x: expected z.M{1,'a'} at U+0073 's' 1,3-3 (3-3)
}

func ExampleX_min_Zero() {
	s := scan.New("  sssome thing")
	s.X(z.M0{' '})
	s.Print()
	s.X(z.M0{'a'})
	s.Print()
	// Output:
	// U+0073 's' 1,3-3 (3-3)
	// U+0073 's' 1,3-3 (3-3)
}

func ExampleX_count() {
	s := scan.New("sssome thing")
	s.X(z.C{3, 's'})
	s.Print()
	s.X(z.C{2, 'o'})
	s.Print()
	// Output:
	// U+006F 'o' 1,4-4 (4-4)
	// scan.x: expected z.C{2,'o'} at U+006F 'o' 1,4-4 (4-4)
}

func ExampleX_new_Line() {
	eol := z.X{z.I{'\n', "\r\n", '\r'}, tk.NL}
	s := scan.New("some\nth\r\ning\rhere")
	s.X("some", eol, "th", eol, "ing", eol, 'h')
	s.Print()
	// Output:
	// U+0065 'e' 4,2-2 (15-15)

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
