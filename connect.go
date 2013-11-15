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

const lineEnding = "\n"

func ConnectServer(w http.ResponseWriter, r *http.Request) {
	addr := fmt.Sprintf("%s:%s", r.FormValue("host"), r.FormValue("port"))
	conn, err := gotelnet.Dial(addr)
	if err != nil {
		log.Println("Error: Dial:", err)
		w.WriteHeader(500)
		return
	}
	websocket.Handler(func(ws *websocket.Conn) {
		makeProxy(ws, conn).run()
	}).ServeHTTP(w, r)
}

type proxy struct {
	client *websocket.Conn
	server gotelnet.Conn
	done   chan bool
}

func makeProxy(client *websocket.Conn, server gotelnet.Conn) *proxy {
	return &proxy{
		client: client,
		server: server,
		done:   make(chan bool),
	}
}

func (pr *proxy) run() {
	go pr.processInput()
	go pr.processOutput()
	<-pr.done
}

func (pr *proxy) processInput() {
	for {
		var msg string
		if rerr := websocket.Message.Receive(pr.client, &msg); rerr != nil {
			log.Println("Error: Receive:", rerr)
			break
		}
		if _, werr := pr.server.Write([]byte(msg + lineEnding)); werr != nil {
			log.Println("Error: Write:", werr)
			break
		}
	}
	pr.done <- true
}

func (pr *proxy) processOutput() {
	r := bufio.NewReader(pr.server)
	for {
		line, _, err := r.ReadLine()
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
		werr := websocket.Message.Send(pr.client, string(runes))
		if werr != nil {
			log.Println("Error: Send:", werr)
			break
		}
	}
	pr.done <- true
}
