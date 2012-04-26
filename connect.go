package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/stesla/gotelnet"
	"io"
	"net/http"
)

func ConnectServer(w http.ResponseWriter, r *http.Request) {
	// TODO: do some sort of argument parsing here
	conn, err := gotelnet.Dial("localhost:2860")
	if err != nil {
		w.WriteHeader(500);
		return
	}
	websocket.Handler(func(ws *websocket.Conn) {
		forward(ws, conn);
	}).ServeHTTP(w, r)
}

func forward(down *websocket.Conn, up gotelnet.Conn) {
	lineEnding := "\n";
	exit := make(chan bool)
	go func() { io.Copy(down, up); exit <- true }()
	go func() {
		for {
			var msg string
			if rerr := websocket.Message.Receive(down, &msg); rerr != nil {
				break
			}
			if _, werr := up.Write([]byte(msg + lineEnding)); werr != nil {
				break
			}
		}
		exit <- true
	}()
	<-exit
}
