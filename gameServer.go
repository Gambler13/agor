package main

import (
	"golang.org/x/image/colornames"
	"image/color"
	"math/rand"
	"time"
)

type GameLoop struct {
	tickRate time.Duration
	quit     chan bool
	World    World
}

type World struct {
	Players  []Player
	CellTree Quadtree
	Food     []EntityImpl
	FoodTree Quadtree
	Bounds   Bounds
}

func InitWorld(b Bounds) World {

	var food []EntityImpl

	i := 0
	for i < 100 {

		var colName color.Color
		if i < 100 {
			colName = colornames.Bisque
		} else if i < 200 {
			colName = colornames.Aliceblue
		} else {
			colName = colornames.Greenyellow
		}

		f := &Food{Entity{
			Id: uint(rand.Uint32()),
			Circle: Circle{
				Radius: 5,
				//Position: getRandomPosition(b, 1),
				Position: Position{X: float64(10.0 * i), Y: float64(10.0 * i)},
			},
			Killer: 0,
			Color:  colName,
		}}

		food = append(food, f)
		i++
	}

	p := Player{
		Id: uint(rand.Int()),
	}

	c := &Cell{
		Entity: Entity{
			Circle: Circle{
				Radius: 50,
				Position: Position{
					X: 100,
					Y: 100,
				},
			},
			Id:     1,
			Killer: 0,
			Color:  colornames.Red,
		},
		Owner: p,
	}

	p.Cells = append(p.Cells, c)

	ct := Quadtree{
		Bounds:     b,
		MaxObjects: 20,
		MaxLevels:  10,
	}

	ft := Quadtree{
		Bounds:     b,
		MaxObjects: 20,
		MaxLevels:  10,
	}

	for i := range food {
		ft.Insert(food[i])
	}

	return World{
		Players:  []Player{p},
		CellTree: ct,
		Food:     food,
		FoodTree: ft,
		Bounds:   b,
	}

}

func (g *GameLoop) run() {
	tickInterval := time.Second / g.tickRate
	timeStart := time.Now().UnixNano()

	ticker := time.NewTicker(tickInterval)

	for {
		select {
		case <-ticker.C:
			now := time.Now().UnixNano()
			// DT in seconds
			delta := float64(now-timeStart) / 1000000000
			timeStart = now
			g.onUpdate(delta)

		case <-g.quit:
			ticker.Stop()
		}
	}
}

func (g *GameLoop) onUpdate(delta float64) {
	p := g.World.Players
	g.World.CellTree.Clear()
	for i := range p {
		for j := range p[i].Cells {
			p[i].Cells[j].move(delta)
			p[i].Cells[j].eat(&g.World.FoodTree)
			//p[i].Cells[j].eat(g.World.CellTree)
			g.World.CellTree.Insert(p[i].Cells[j])
		}
	}

}

type Player struct {
	Id    uint
	Mouse Position
	Cells []*Cell
}

func (p *Player) getCenter() Position {
	pos := make([]Position, len(p.Cells))
	for i := range p.Cells {
		pos[i] = p.Cells[i].Position
	}

	return centroid(pos)

}
func (p *Player) getDistance(pos Position) float64 {
	return getDistance(p.Mouse, pos)
}
