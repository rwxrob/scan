package pg_test

import (
	"github.com/rwxrob/scan"
	"github.com/rwxrob/scan/pg"
)

func ExampleUGraphic() {
	s := scan.New(`some thing`)
	s.X(pg.UGraphic)
	s.Print()
	// Output:
	// U+006F 'o' 1,2-2 (2-2)
}
