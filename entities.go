package main

import (
	"fmt"
	"image/color"
)

type Entity struct {
	Circle
	Id     uint
	Killer uint
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
	Owner Player
}

func (c *Cell) onAdd() {
	panic("implement me")
}

func (c *Cell) onRemove() {
	panic("implement me")
}

func (c *Cell) onConsume(entity *Entity) {
	c.Killer = entity.Id
}

func (c *Cell) move(delta float64) {
	movement := 1.0 / c.Radius * delta * 150
	c.Position.X += movement
	//c.Position.Y += movement
}

func (c *Cell) eat(qt *Quadtree) {
	intersections := qt.RetrieveIntersections(c)
	for i := range intersections {
		e := intersections[i]
		ee := e.getEntity()
		if ee.Radius < c.Radius && ee.Killer == 0 {
			c.Radius += ee.Radius * 0.075
			e.onConsume(c.getEntity())
			fmt.Printf("eat eat %d\n", ee.Id)

		}
	}

}

type Food struct {
	Entity
}

func (f *Food) onAdd() {
	panic("implement me")
}

func (f *Food) onRemove() {
	panic("implement me")
}

func (f *Food) onConsume(entity *Entity) {
	f.Killer = entity.Id
}

func (f *Food) onUpdate(qt *Quadtree, delta float64) {
}

type Virus struct {
	Entity
}

func (v *Virus) onAdd() {
	panic("implement me")
}

func (v *Virus) onRemove() {
	panic("implement me")
}

func (v *Virus) onConsume(entity *Entity) {
	panic("implement me")
}

type EntityImpl interface {
	onAdd()
	onRemove()
	onConsume(entity *Entity)
	getEntity() *Entity
}
