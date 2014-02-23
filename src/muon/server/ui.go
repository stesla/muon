package main

import (
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"io"
	"log"
	"muon/model"
	"net/http"
	"os"
	"path"
)

type uiKey int

const (
	uiUI uiKey = 0 + iota
)

type UI struct {
	dir string
	router   *mux.Router
}

func NewUI(dir string) http.Handler {
	r := mux.NewRouter()

	r.Path("/login").Name("login").HandlerFunc(login)

	return &UI{dir, r}
}

func (ui *UI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context.Set(r, uiUI, ui)

	ui.router.ServeHTTP(w, r)
}

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		dir := context.Get(r, uiUI).(*UI).dir
		file, err := os.Open(path.Join(dir, "login.html"))
		if err != nil {
			log.Println("error opening file:", err)
		}
		io.Copy(w, file)

	case "POST":
		login := r.FormValue("login")
		pw := r.FormValue("pw")
		user, err := model.AuthUser(login, pw)
		if err != nil {
			fmt.Fprintln(w, err)
		} else {
			fmt.Fprintln(w, user.Login, user.Email)
		}

	default:
		w.Header().Set("Allow", "GET, POST")
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
