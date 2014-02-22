package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"muon/model"
	"net/http"
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

	r.Methods("GET").Path("/login").
		Name("showLogin").HandlerFunc(showLogin)

	r.Methods("POST").Path("/login").
		Name("authLogin").HandlerFunc(authLogin)

	return &UI{dir, r}
}

func (ui *UI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context.Set(r, uiUI, ui)

	ui.router.ServeHTTP(w, r)
}

func authLogin(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	pw := r.FormValue("pw")
	user, err := model.AuthUser(login, pw)
	if err != nil {
		fmt.Fprintln(w, err)
	} else {
		fmt.Fprintln(w, user.Login, user.Email)
	}
}

func showLogin(w http.ResponseWriter, r *http.Request) {
	dir := context.Get(r, uiUI).(*UI).dir
	file := path.Join(dir, "login.html")
	tpl, err := template.ParseFiles(file)
	if err != nil {
		log.Println("error loading template:", err)
		w.WriteHeader(500)
		return
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, nil); err != nil {
		log.Println("error rendering template:", err)
		w.WriteHeader(500)
		return
	} else {
		buf.WriteTo(w)
	}
}
