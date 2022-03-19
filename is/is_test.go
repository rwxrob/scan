package z_test

import (
	"fmt"

	z "github.com/rwxrob/scan/is"
)

// P X A Y N I O T Ti R MM M M0 M1 C

func ExampleP_String() {
	fmt.Printf("%v\n", z.P{2, 's'})
	// Output:
	// z.P{2,'s'}
}

func ExampleX_String() {
	fmt.Printf("%v\n", z.X{'g', 'o'})
	// Output:
	// z.X{'g','o'}
}

func ExampleA_String() {
	fmt.Printf("%v\n", z.A{5})
	// Output:
	// z.A{5}
}

func ExampleY_String() {
	fmt.Printf("%v\n", z.Y{'s', "foo"})
	// Output:
	// z.Y{'s',"foo"}
}

func ExampleN_String() {
	fmt.Printf("%v\n", z.N{'s', "foo"})
	// Output:
	// z.N{'s',"foo"}
}

func ExampleI_String() {
	fmt.Printf("%v\n", z.I{'s', "foo"})
	// Output:
	// z.I{'s',"foo"}
}

func ExampleO_String() {
	fmt.Printf("%v\n", z.O{'s', "foo"})
	// Output:
	// z.O{'s',"foo"}
}

func ExampleT_String() {
	fmt.Printf("%v\n", z.T{'s'})
	// Output:
	// z.T{'s'}
}

func ExampleTi_String() {
	fmt.Printf("%v\n", z.Ti{'s'})
	// Output:
	// z.Ti{'s'}
}

func ExampleR_String() {
	fmt.Printf("%v\n", z.R{'l', 'o'})
	// Output:
	// z.R{'l','o'}
}

func ExampleMM_String() {
	fmt.Printf("%v\n", z.MM{1, 3, 's'})
	// Output:
	// z.MM{1,3,'s'}
}

func ExampleM_String() {
	fmt.Printf("%v\n", z.M{1, 's'})
	// Output:
	// z.M{1,'s'}
}

func ExampleM0_String() {
	fmt.Printf("%v\n", z.M0{'s'})
	// Output:
	// z.M0{'s'}
}

func ExampleM1_String() {
	fmt.Printf("%v\n", z.M1{'s'})
	// Output:
	// z.M1{'s'}
}

func ExampleC_String() {
	fmt.Printf("%v\n", z.C{2, 's'})
	// Output:
	// z.C{2,'s'}
}
