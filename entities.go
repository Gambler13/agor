package main

import (
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"math"
	"math/rand"
)

type Entity struct {
	Circle
	Id     int
	Killer int
	color  color.Color
}

func (e *Entity) getEntity() *Entity {
	return e
}

func (e Entity) Bounds() image.Rectangle {

	xMin := float64(e.X) - e.Radius
	yMin := float64(e.Y) - e.Radius
	xMax := float64(e.X) + e.Radius
	yMax := float64(e.Y) + e.Radius
	return image.Rectangle{
		Min: image.Point{
			X: int(xMin),
			Y: int(yMin),
		},
		Max: image.Point{
			X: int(xMax),
			Y: int(yMax),
		},
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
			color:  colornames.Beige,
		},
		Owner: owner,
	}
}

func (c *Cell) onConsume(entity *Entity) {
	c.Killer = entity.Id
}

func (c *Cell) move(delta float64) {

	movement := 1.0 / c.Radius * delta * 3500
	mP := c.Owner.Mouse

	vec := math.Sqrt(q(mP.X) + q(mP.Y))

	if vec != 0.0 {
		x := movement * (float64(mP.X) / vec)
		y := movement * (float64(mP.Y) / vec)
		c.Position.X += int(x)
		c.Position.Y += int(y)
	}

}

func (c *Cell) eat(qt *QuadTree) {
	intersections := qt.query(c.Bounds())
	for i := range intersections {
		e := intersections[i]
		//Check if cell eat itself
		interC, ok := e.(*Cell)
		if ok {
			if interC.Owner.Id == c.Owner.Id {
				continue
			}
		}
		entImpl := e.getEntity()

		if entImpl.Radius < c.Radius && entImpl.Killer == 0 {
			c.Owner.distributeFood(e)
			c.Owner.foodEaten++
			e.onConsume(c.getEntity())
		}
	}
}

func (c *Cell) rectangle() image.Rectangle {
	xMin := float64(c.X) - c.Radius
	yMin := float64(c.Y) - c.Radius
	xMax := float64(c.X) + c.Radius
	yMax := float64(c.Y) + c.Radius
	return image.Rectangle{
		Min: image.Point{
			X: int(xMin),
			Y: int(yMin),
		},
		Max: image.Point{
			X: int(xMax),
			Y: int(yMax),
		},
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
			color:  c.color,
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

func (f *Food) onUpdate(qt *QuadTree, delta float64) {
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
