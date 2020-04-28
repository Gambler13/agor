package api

type Entity struct {
	X      float64
	Y      float64
	Radius float64
	Color  string
}

type Mouse struct {
	X float64
	Y float64
}

type GameStats struct {
	PlayerId   string
	Mass       float64
	FoodEaten  int
	CellsEaten int
	Rank       int
}
