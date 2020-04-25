package main

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
)

var canvasWidth = 640.0
var canvasHeight = 400.0

func main() {

	w := InitWorld(Bounds{
		X:      0,
		Y:      0,
		Width:  1200,
		Height: 900,
	})
	gl := GameLoop{
		tickRate: 50,
		quit:     nil,
		World:    w,
	}

	go gl.run()

	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     int(canvasWidth),
		Height:    int(canvasHeight),
		FrameRate: 30,
		Title:     "Hello Canvas!",
	})

	c.Setup(func(ctx *canvas.Context) {
		ctx.SetColor(colornames.White)
		ctx.Clear()
		ctx.SetColor(colornames.Green)
		ctx.SetLineWidth(1)
	})

	c.Draw(func(ctx *canvas.Context) {
		ctx.Clear()

		p := w.Players[0].getCenter()
		dx := canvasWidth/2 - p.X
		dy := canvasHeight/2 - p.Y
		width := 1200.0
		height := 800.0

		ctx.Push()
		ctx.Translate(dx, dy)
		ctx.SetColor(colornames.White)
		ctx.DrawRectangle(0, 0, width, height)
		ctx.Fill()
		ctx.Pop()

		ctx.Push()
		ctx.Translate(dx, dy)
		ctx.SetColor(colornames.Blue)
		ctx.DrawCircle(p.X, p.Y, 3)
		ctx.Stroke()
		ctx.DrawCircle(100, 100, 2)
		ctx.Stroke()
		ctx.Pop()

		ctx.Push()
		ents := gl.World.CellTree.RetrieveViewIntersections(Bounds{
			X:      p.X,
			Y:      p.Y,
			Width:  300,
			Height: 300,
		})
		ents2 := gl.World.FoodTree.RetrieveViewIntersections(Bounds{
			X:      p.X,
			Y:      p.Y,
			Width:  300,
			Height: 300,
		})

		ents = append(ents, ents2...)

		ctx.Push()
		ctx.Translate(dx, dy)
		ctx.SetColor(colornames.Orange)
		ctx.DrawRectangle(p.X-150, p.Y-150, 300, 300)
		ctx.Stroke()
		ctx.Pop()

		for _, e := range ents {
			entity := e.getEntity()
			ctx.SetColor(entity.Color)

			pos := entity.Position

			ctx.Push()
			ctx.Translate(dx, dy)
			ctx.DrawCircle(pos.X, pos.Y, e.getEntity().Radius)
			ctx.Fill()

			if e.getEntity().Killer != 0 {
				ctx.SetColor(colornames.Black)
				ctx.DrawCircle(pos.X, pos.Y, e.getEntity().Radius)
				ctx.Stroke()
			}
			ctx.Pop()

		}

		ctx.Pop()

		if ctx.IsKeyPressed(pixelgl.KeyUp) {
			ctx.Push()
			ctx.SetColor(colornames.White)
			ctx.Clear()
			ctx.Pop()
		}

	})
}
