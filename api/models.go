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
	PlayerID   string
	NumPlayers int
	Mass       float64
	FoodEaten  int
	CellsEaten int
	Rank       int
}

type MouseEvent struct {
	SeqID uint32
	X     float32
	Y     float32
}
