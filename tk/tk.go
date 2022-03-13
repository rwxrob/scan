package tk

const (

	// EOD is a special value that is returned when the end of data is
	// reached enabling functional parser functions to look for it reliably
	// no matter what is being parsed. Since rune is alias for int32 and
	// Unicode (currently) ends at \U+FFFD we are safe to use the largest
	// possible valid rune value.
	EOD rune = 1<<31 - 1 // max int32
	ANY                  // represents any valid rune
)
