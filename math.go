package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
)

type Position struct {
	X int
	Y int
}

type Position64 struct {
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

func q(a float64) float64 {
	return a * a
}

func getRandomPosition(r image.Rectangle, padding float64) Position {
	max := r.Max
	min := r.Min
	x := rand.Intn(max.X-min.X+1) + min.X
	y := rand.Intn(max.Y-min.Y+1) + min.Y

	return Position{
		X: x,
		Y: y,
	}
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

	var totalPoints = len(points)
	center.X = center.X / totalPoints
	center.Y = center.Y / totalPoints

	return center
}

func hexColor(c color.Color) string {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("#%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B)
}

func randomColor() color.Color {
	r := uint8(rand.Intn(255))
	g := uint8(rand.Intn(255))
	b := uint8(rand.Intn(255))
	return color.RGBA{
		R: r,
		G: g,
		B: b,
		A: 1,
	}
}
