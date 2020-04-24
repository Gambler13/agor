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
	Quadtree Quadtree
	Bounds   Bounds
}

func InitWorld(b Bounds) World {

	var e []EntityImpl

	i := 0
	for i < 350 {

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
				Position: Position{X: float64(10.0*i) - 250, Y: float64(10.0*i) - 250},
			},
			Killer: 0,
			Color:  colName,
		}}

		e = append(e, f)
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
			Id:     uint(rand.Uint32()),
			Killer: 0,
			Color:  colornames.Red,
		},
		Owner: p,
	}

	p.Cells = append(p.Cells, c)

	e = append(e, c)

	qt := Quadtree{
		Bounds:     b,
		MaxObjects: 20,
		MaxLevels:  10,
		Objects:    e,
	}

	return World{
		Quadtree: qt,
		Bounds:   b,
		Players:  []Player{p},
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

	qt := g.World.Quadtree

	g.updateQTree(qt, delta)
}

func (g *GameLoop) updateQTree(qt Quadtree, delta float64) {

	for i := range qt.Objects {
		if qt.Objects[i].getEntity().Killer == 0 {
			qt.Objects[i].onUpdate(&qt, delta)
		} else {
			copy(qt.Objects[i:], qt.Objects[i+1:])
		}
	}

	if len(qt.Nodes) > 0 {
		for i := 0; i < len(qt.Nodes); i++ {
			g.updateQTree(qt.Nodes[i], delta)
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
