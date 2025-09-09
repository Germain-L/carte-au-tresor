package models

// Direction constants for better performance than string comparisons
type Direction int

const (
	North Direction = iota
	East
	South
	West
)

// DirectionStrings maps Direction constants to their string representations
var DirectionStrings = map[Direction]string{
	North: "N",
	East:  "E", 
	South: "S",
	West:  "O",
}

// StringToDirection maps string directions to Direction constants
var StringToDirection = map[string]Direction{
	"N": North,
	"E": East,
	"S": South,
	"O": West,
}

// String returns the string representation of a Direction
func (d Direction) String() string {
	return DirectionStrings[d]
}