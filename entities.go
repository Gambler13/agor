package main

import (
	"bytes"
	"encoding/binary"
	"image"
	"math"
	"math/rand"
)

type Entity struct {
	Circle
	Id     int
	Killer int
	color  byte //Real value doesn't matter
}

func (e *Entity) getEntity() *Entity {
	return e
}

//Return Position castet to int
func (e *Entity) PosInt() (x, y int) {
	return int(e.X), int(e.Y)
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
			color:  randomLutIndex(),
		},
		Owner: owner,
	}
}

func (c *Cell) onConsume(entity *Entity) {
	c.Killer = entity.Id
}

func (c *Cell) move(delta float64) {

	movement := 1.0 / c.Radius * delta * 1500
	mP := c.Owner.Mouse

	vec := math.Sqrt(q(mP.X) + q(mP.Y))

	if vec != 0.0 {
		x := movement * (mP.X / vec)
		y := movement * (mP.Y / vec)
		c.Position.X += x
		c.Position.Y += y
	}

}

//32b X
//32b Y
//64b Radius
//8b  Color
//----
//17 byte
func (e *Entity) getByteSize() int {
	return 17
}

func (e *Entity) getByte() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, e.getByteSize()))

	binary.Write(buf, binary.BigEndian, uint32(e.X))
	binary.Write(buf, binary.BigEndian, uint32(e.Y))
	binary.Write(buf, binary.BigEndian, e.Radius)
	binary.Write(buf, binary.BigEndian, e.color)

	return buf.Bytes()
}

func (c *Cell) eat(qt *QuadTree) {
	intersections := qt.query(c.Bounds())
	for i := range intersections {
		e := intersections[i]
		victim, ok := e.(*Cell)
		if ok {
			//Check if cell eat itself
			if victim.Owner.Id == c.Owner.Id {
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
	xMin := c.X - c.Radius
	yMin := c.Y - c.Radius
	xMax := c.X + c.Radius
	yMax := c.Y + c.Radius
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
	getByteSize() int
	getByte() []byte
}
