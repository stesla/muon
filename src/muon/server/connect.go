package main

import (
	"bufio"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/stesla/gotelnet"
	"io"
	"log"
	"net/http"
	"unicode/utf8"
)

//TODO: remove error logging that shouldn't really be there

const lineEnding = "\n"

func connect(w http.ResponseWriter, r *http.Request) {
	addr := fmt.Sprintf("%s:%s", r.FormValue("host"), r.FormValue("port"))
	conn, err := gotelnet.Dial(addr)
	if err != nil {
		log.Println("Error: Dial:", err)
		w.WriteHeader(500)
		return
	}
	websocket.Handler(func(ws *websocket.Conn) {
		makeProxy(wsconn{ws}, conn).run()
	}).ServeHTTP(w, r)
}

type wsconn struct {
	*websocket.Conn
}

func (ws wsconn) Send(in string) error {
	return websocket.Message.Send(ws.Conn, in)
}

func (ws wsconn) Receive(out *string) error {
	return websocket.Message.Receive(ws.Conn, out)
}

type ReceiveSender interface {
	Receive(out *string) error
	Send(in string) error
}

type proxy struct {
	client ReceiveSender
	server io.ReadWriter
	done   chan bool
}

func makeProxy(client ReceiveSender, server io.ReadWriter) *proxy {
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
		if rerr := pr.client.Receive(&msg); rerr != nil {
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
		werr := pr.client.Send(string(runes))
		if werr != nil {
			log.Println("Error: Send:", werr)
			break
		}
	}
	pr.done <- true
}
