package models

type Mountain struct {
	Row, Col int
}

type Treasure struct {
	Row, Col int
	Count    int
}

type Adventurer struct {
	Name      string
	Row, Col  int
	Direction string
	Actions   string
	Treasures []Treasure
}

type MapCell struct {
	Mountain   *Mountain
	Treasure   *Treasure
	Adventurer *Adventurer
}

type MapData struct {
	Width, Height int
	Mountains     []Mountain
	Treasures     []Treasure
	Adventurers   []Adventurer
}
