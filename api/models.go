package api

type Entity struct {
	X      int
	Y      int
	Radius float64
	Color  string
}

type Mouse struct {
	X float64
	Y float64
}

type GameStats struct {
	PlayerId   string
	NumPlayers int
	Mass       float64
	FoodEaten  int
	CellsEaten int
	Rank       int
}

type MouseEvent struct {
	X float32
	Y float32
}
