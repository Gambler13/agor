package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
)

type Position struct {
	X float64
	Y float64
}

type Circle struct {
	Position
	Radius float64
}

type Rectangle struct {
	Position
	Width  float64
	Height float64
}

func getDistance(a, b Position) float64 {
	return math.Sqrt(q(b.X-a.X) * q(b.Y-a.Y))
}

func sub(a, b Position) Position {
	return Position{
		X: a.X - b.X,
		Y: a.Y - b.Y,
	}
}

func q(a float64) float64 {
	return a * a
}

func getRandomPosition(bounds Bounds, padding float64) Position {
	maxX := int(bounds.Width - padding)
	minX := int(padding)
	maxY := int(bounds.Height - padding)
	minY := int(padding)
	x := float64(rand.Intn(maxX-minX+1) + minX)
	y := float64(rand.Intn(maxY-minY+1) + minY)

	return Position{
		X: x,
		Y: y,
	}
}

func intersects(c Circle, r Rectangle) bool {

	cDistanceX := math.Abs(c.X - r.X)
	cDistanceY := math.Abs(c.Y - r.Y)

	if cDistanceX > (r.Width/2 + c.Radius) {
		return false
	}

	if cDistanceY > (r.Height/2 + c.Radius) {
		return false
	}

	if cDistanceX <= (r.Width / 2) {
		return true
	}

	if cDistanceY <= (r.Height / 2) {
		return true
	}

	cornerDistanceSq := q(cDistanceX-r.Width/2) + q(cDistanceY-r.Height/2)

	return cornerDistanceSq <= q(c.Radius)
}

func intersectsPoint(p Position, r Bounds) bool {
	if p.X < r.X {
		return false
	}

	if p.X > r.X+r.Width {
		return false
	}

	if p.Y < r.Y {
		return false
	}

	if p.Y > r.Y+r.Width {
		return false
	}

	return true

}

func centroid(points []Position) Position {
	if len(points) == 1 {
		return points[0]
	}

	var center Position

	for i := 0; i < len(points); i += 1 {
		center.X += points[i].X
		center.Y += points[i].Y
	}

	var totalPoints = float64(len(points))
	center.X = center.X / totalPoints
	center.Y = center.Y / totalPoints

	return center
}

func hexColor(c color.Color) string {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("#%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B)
}
