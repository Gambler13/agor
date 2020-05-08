package main

import (
	"fmt"
	"github.com/gambler13/agor/conf"
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

const AgorConfigEnv = "AGOR_CONFIG"
const DefaultConfig = "/etc/agor/config/config.yaml"

var canvasWidth = 640.0
var canvasHeight = 400.0

func main() {
	Log = logrus.New()

	configFile := os.Getenv(AgorConfigEnv)
	if configFile == "" {
		Log.Warnf("no config file provided with %s env variable, use default config %s", AgorConfigEnv, DefaultConfig)
		configFile = DefaultConfig
	}

	conf, err := conf.Load(configFile)
	if err != nil {
		panic(fmt.Errorf("could not load config file: %w", err))
	}

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
