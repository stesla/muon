package main

import (
	"bufio"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/stesla/gotelnet"
	"log"
	"net/http"
	"unicode/utf8"
)

func ConnectServer(w http.ResponseWriter, r *http.Request) {
	addr := fmt.Sprintf("%s:%s", r.FormValue("host"), r.FormValue("port"))
	conn, err := gotelnet.Dial(addr)
	if err != nil {
		log.Println("Error: Dial:", err)
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
	go func() {
		b := bufio.NewReader(up)
		for {
			line, _, err := b.ReadLine()
			if err != nil {
				log.Println("Error: ReadLine:", err)
				break
			}
			runes := []rune{}
			for len(line) > 0 {
				r, n := utf8.DecodeRune(line)
				if r == utf8.RuneError {
					r = '?'
				}
				runes = append(runes, r)
				line = line[n:]
			}
			werr := websocket.Message.Send(down, string(runes))
			if werr != nil {
				log.Println("Error: Send:", werr)
				break
			}
		}
		exit <- true
	}()
	go func() {
		for {
			var msg string
			if rerr := websocket.Message.Receive(down, &msg); rerr != nil {
				log.Println("Error: Receive:", rerr);
				break
			}
			if _, werr := up.Write([]byte(msg + lineEnding)); werr != nil {
				log.Println("Error: Write:", werr);
				break
			}
		}
		exit <- true
	}()
	<-exit
}
