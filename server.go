package main

import (
	"encoding/json"
	"fmt"
	"hexhibit.xyz/agor/api"
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
)

func startServer(w *World) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("connected:", s.ID())

		w.addNewPlayer(s)

		return nil
	})

	server.OnEvent("/", "position", func(s socketio.Conn, msg string) {
		var pos api.Mouse
		json.Unmarshal([]byte(msg), &pos)
		w.handlePosition(s.ID(), Position{
			X: pos.X,
			Y: pos.Y,
		})

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
		s.Close()
		return last
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
	go server.Serve()
	defer server.Close()

	http.Handle("/socket/", corsMiddleware(server))
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8008...")
	log.Fatal(http.ListenAndServe(":8008", nil))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8084")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

		next.ServeHTTP(w, r)
	})
}
