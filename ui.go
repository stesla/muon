package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
)

type UIServer struct {}

func NewUIServer() http.Handler {
	return &UIServer{}
}

func (s *UIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	dir := path.Join(dataDir, "templates")
	t, err := template.ParseFiles(path.Join(dir, "index.html"))
	if err != nil {
		log.Println("ParseFiles:", err)
		http.Error(w, "Template Error", http.StatusInternalServerError)
	}
	t.Execute(w, nil)
}
