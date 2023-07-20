package util

type Direction int

// directions iterate clockwise from north
const (
	North Direction = 0
	West  Direction = 1
	South Direction = 2
	East  Direction = 3
)

// Directions
// can use this to iterate clearly over each cardinal direction (i mean or just use a 'for < 4' loop)
var Directions = []Direction{
	North,
	West,
	South,
	East,
}
