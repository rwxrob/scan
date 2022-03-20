package scan_test

import (
	"fmt"

	"github.com/rwxrob/scan"
	z "github.com/rwxrob/scan/is"
)

func ExampleX_parse_Single() {
	s := scan.New("some thing")
	s.X(z.P{1, z.X{'s', 'o', "me"}})
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

func ExampleX_parse_Nested_in_Sequence() {
	s := scan.New("some thing")
	s.X(z.P{2, z.X{"some", ' ', z.P{3, "th"}}})
	s.Print()
	fmt.Println(s.Tree.Root.Count)
	s.Tree.Root.Print()
	// Output:
	// U+0069 'i' 1,8-8 (8-8)
	// 1
	// {"T":1,"N":[{"T":2,"V":"some th","N":[{"T":3,"V":"th"}]}]}
}

func ExampleX_parse_Nested_Two() {
	s := scan.New("some thing")
	s.X(z.P{2, z.X{z.P{3, "some"}, ' ', z.P{3, "th"}}})
	s.Print()
	fmt.Println(s.Tree.Root.Count)
	s.Tree.Root.Print()
	// Output:
	// U+0069 'i' 1,8-8 (8-8)
	// 1
	// {"T":1,"N":[{"T":2,"V":"some th","N":[{"T":3,"V":"some"},{"T":3,"V":"th"}]}]}
}

func ExampleX_parse_Failure() {

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
	// scan.x: expected z.I{' ','\t','\r','\n'} at U+0045 'E' 1,3-3 (3-3)
	// {"T":1}
}

func ExampleX_parse_Nested_EOD() {

	const WORD = 2
	const CHAR = 3

	ch := z.P{CHAR, z.R{'a', 'z'}}
	word := z.P{WORD, z.M1{ch}}
	ws := z.I{' ', '\t', '\r', '\n'}

	s := scan.New("me")
	s.X(z.X{word, ws}, word)
	//s.X(z.X{word, ws}, word)
	s.Print()
	s.Tree.Root.Print()

	// Output:
	// scan.x: expected z.I{' ','\t','\r','\n'} at <EOD>
	// {"T":1}
}

func ExampleX_parse_Nested_Complex() {

	const WORD = 2
	const CHAR = 3

	ch := z.P{CHAR, z.R{'a', 'z'}}
	word := z.P{WORD, z.M1{ch}}
	ws := z.I{' ', '\t', '\r', '\n'}

	s := scan.New("go me again")
	s.X(z.X{word, ws})
	s.Print()
	//nodes := s.Tree.Root.Nodes()
	//fmt.Println(s.Tree.Root.Count, nodes[0].Count, nodes[1].Count)
	s.Tree.Root.Print()

	// Output:
	// U+0020 ' ' 1,6-6 (6-6)
	// 2 2 2
	// {"T":1,"N":[{"T":2,"V":"go","N":[{"T":3,"V":"g"},{"T":3,"V":"o"}]},{"T":2,"V":"me","N":[{"T":3,"V":"m"},{"T":3,"V":"e"}]}]}
}

/*

func ExampleX_parse_Nested_Complex_Revert() {

	const WORD = 2
	const CHAR = 3

	ch := z.P{CHAR, z.R{'a', 'z'}}
	word := z.P{WORD, z.M1{ch}}
	ws := z.I{' ', '\t', '\r', '\n'}

	s := scan.New("go me ok doke")
	s.X(z.M0{"go"}, ' ', z.M0{z.X{word, ws}})
	s.Print()
	s.Tree.Root.Print()

	// Output:
	// U+0064 'd' 1,10-10 (10-10)
	// {"T":1,"N":[{"T":2,"V":"me","N":[{"T":3,"V":"m"},{"T":3,"V":"e"}]},{"T":2,"V":"ok","N":[{"T":3,"V":"o"},{"T":3,"V":"k"}]}]}

}

func ExampleX_parse_Min() {
	const SINGLE = 2
	s := scan.New("g m doke")
	s.X(z.M0{z.X{z.P{SINGLE, z.R{'a', 'z'}}, ' '}})
	s.Print()
	s.Tree.Root.Print()
	// Output:
	// U+0064 'd' 1,5-5 (5-5)
	// {"T":1,"N":[{"T":2,"V":"g"},{"T":2,"V":"m"}]}

}

func ExampleX_parse_Nested_Complex_Multipart() {

	const WORD = 2
	const CHAR = 3

	ch := z.P{CHAR, z.R{'a', 'z'}}
	word := z.P{WORD, z.M1{ch}}
	ws := z.I{' ', '\t', '\r', '\n'}
	phrase := z.X{z.M0{z.X{word, ws}}, word}

	s := scan.New("go me again")
	s.X(phrase)
	s.Print()
	//nodes := s.Tree.Root.Nodes()
	//fmt.Println(s.Tree.Root.Count, nodes[0].Count, nodes[1].Count)
	s.Tree.Root.Print()

	// Output:
	// U+0020 ' ' 1,6-6 (6-6)
	// 2 2 2
	// {"T":1,"N":[{"T":2,"V":"go","N":[{"T":3,"V":"g"},{"T":3,"V":"o"}]},{"T":2,"V":"me","N":[{"T":3,"V":"m"},{"T":3,"V":"e"}]}]}
}
*/
