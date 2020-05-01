package main

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
	"image"
)

var Log *logrus.Logger

var canvasWidth = 640.0
var canvasHeight = 400.0

func main() {

	Log = logrus.New()

	bounds := image.Rectangle{
		Min: image.Point{},
		Max: image.Point{
			X: 1200,
			Y: 800,
		},
	}

	w := InitWorld(bounds)
	gl := GameLoop{
		tickRate:       50,
		quit:           nil,
		World:          w,
		PositionCh:     make(chan PositionMsg, 300),
		AddPlayerCh:    make(chan socketio.Conn, 20),
		RemovePlayerCh: make(chan string, 0),
	}

	go gl.run()

	startServer(gl.AddPlayerCh, gl.RemovePlayerCh, gl.PositionCh)

}
