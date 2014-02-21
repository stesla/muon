package main

import (
	"net/http"
	"strings"
)

type StaticFileHandler string

func (h StaticFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if upath := r.URL.Path; !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	f, err := http.Dir(h).Open(r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()

	d, err1 := f.Stat()
	if err1 != nil {
		http.NotFound(w, r)
		return
	}

	if d.IsDir() {
		http.NotFound(w, r)
		return
	}

	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}
