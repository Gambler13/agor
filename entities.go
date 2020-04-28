package main

import (
	"golang.org/x/image/colornames"
	"image/color"
	"math"
	"math/rand"
)

type Entity struct {
	Circle
	Id     int
	Killer int
	Color  color.Color
}

func (e *Entity) getEntity() *Entity {
	return e
}

func (e Entity) Bounds() Bounds {
	return Bounds{
		X:      e.Position.X - e.Radius,
		Y:      e.Position.Y - e.Radius,
		Width:  e.Radius * 2,
		Height: e.Radius * 2,
	}
}

type Cell struct {
	Entity
	Owner *Player
}

func (w *World) NewCell(owner *Player) Cell {
	return Cell{
		Entity: Entity{
			Circle: Circle{
				Radius:   20,
				Position: getRandomPosition(w.Bounds, 10),
			},
			Id:     rand.Int(),
			Killer: 0,
			Color:  colornames.Beige,
		},
		Owner: owner,
	}
}

func (c *Cell) onConsume(entity *Entity) {
	c.Killer = entity.Id
}

func (c *Cell) move(delta float64) {

	movement := 1.0 / c.Radius * delta * 350
	mP := c.Owner.Mouse

	vec := math.Sqrt(q(mP.X) + q(mP.Y))

	if vec != 0.0 {
		x := movement * (mP.X / vec)
		y := movement * (mP.Y / vec)
		c.Position.X += x
		c.Position.Y += y
	}

}

func (c *Cell) eat(qt *Quadtree) {
	intersections := qt.RetrieveIntersections(c)
	for i := range intersections {
		e := intersections[i]
		entImpl := e.getEntity()

		//Check if cell eat itself
		interC, ok := e.(*Cell)
		if ok {
			if interC.Owner.Id == c.Owner.Id {
				continue
			}
		}

		if entImpl.Radius < c.Radius && entImpl.Killer == 0 {
			c.Owner.distributeFood(e)
			e.onConsume(c.getEntity())
		}
	}
}

func (c *Cell) split() *Cell {
	newRad := c.Radius / 2
	c.Radius = newRad

	newCell := &Cell{
		Entity: Entity{
			Circle: Circle{
				Position: Position{
					X: c.X,
					Y: c.Y,
				},
				Radius: newRad,
			},
			Id:     rand.Int(),
			Killer: 0,
			Color:  c.Color,
		},
		Owner: c.Owner,
	}

	return newCell
}

type Food struct {
	Entity
}

func (f *Food) onConsume(entity *Entity) {
	f.Killer = entity.Id
}

func (f *Food) onUpdate(qt *Quadtree, delta float64) {
}

type Virus struct {
	Entity
}

func (v *Virus) onConsume(entity *Entity) {
	panic("implement me")
}

type EntityImpl interface {
	onConsume(entity *Entity)
	getEntity() *Entity
}
