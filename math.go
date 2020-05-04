package main

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"strings"
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

var LUT = map[byte]string{
	0: "#ff99f0",
	1: "#37c9f4",
	2: "#f6e327",
	3: "#ad3bee",
	4: "#f22727",
}

func toHexColor(c color.Color) string {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("#%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B)
}

func fromHexColor(s string) color.RGBA {
	s = strings.ReplaceAll(s, "#", "")
	b, _ := hex.DecodeString(s)
	return color.RGBA{b[0], b[1], b[2], b[3]}
}

func randomLutIndex() byte {
	return byte(rand.Intn(len(LUT)))
}
func randomLutColor() color.RGBA {
	n := rand.Intn(4)
	return fromHexColor(LUT[byte(n)])
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
