package main

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

var canvasWidth = 640.0
var canvasHeight = 400.0

func main() {

	conf := Config{
		Port:           8008,
		AllowedOrigins: []string{},
		World: WorldConfig{
			Width:     1200,
			Height:    800,
			MaxPlayer: 100,
			Food:      500,
		},
	}

	Log = logrus.New()

	w := InitWorld(conf)
	gl := GameLoop{
		tickRate:       50,
		quit:           nil,
		World:          w,
		PositionCh:     make(chan PositionMsg, 300),
		AddPlayerCh:    make(chan socketio.Conn, 20),
		RemovePlayerCh: make(chan string, 0),
	}

	go gl.run()

	startServer(gl.AddPlayerCh, gl.RemovePlayerCh, gl.PositionCh, conf.Port)

}
