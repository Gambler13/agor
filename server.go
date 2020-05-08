package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/googollee/go-socket.io/parser"
	"log"
	"net/http"
	"os"

	"github.com/googollee/go-socket.io"
)

func startServer(addPlayer chan socketio.Conn, removePlayer chan string, position chan PositionMsg, port int) {
	Log.Info("start server")
	server, err := socketio.NewServer(nil)
	if err != nil {
		Log.Fatalf("could no create socket.io server: %v", err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		Log.Debugf("new connection:", s.ID())
		addPlayer <- s
		return nil
	})

	server.OnEvent("/", "position", func(s socketio.Conn, msg parser.Buffer) {

		a := struct {
			X float32
			Y float32
		}{}

		buf := bytes.NewReader(msg.Data)
		err := binary.Read(buf, binary.LittleEndian, &a)
		if err != nil {

		}

		pmsg := PositionMsg{PlayerID: s.ID(),
			X: float64(a.X),
			Y: float64(a.Y),
		}

		position <- pmsg

	})

	server.OnEvent("/", "split", func(s socketio.Conn, msg string) string {
		return ""
	})

	server.OnEvent("/", "diet", func(s socketio.Conn, msg string) string {
		return ""
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		removePlayer <- s.ID()
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		Log.Warnf("socket.io error: %v", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		Log.Infof("socket closed connection, ID: %s, reason: %s", s.ID(), reason)
		removePlayer <- s.ID()
		s.Close()
	})
	go server.Serve()
	defer server.Close()

	http.Handle("/socket/", corsMiddleware(server))
	http.Handle("/", http.FileServer(http.Dir("./assets")))
	Log.Infof("listen and serve on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("AGOR_ORIGIN"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
