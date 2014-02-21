package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
)

var (
	listenAddr = flag.String("http", ":8080", "")
	templatesDir    = flag.String("templates", "./templates", "")
	assetsDir  = flag.String("assets", "./assets", "")
)

func Usage() {
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if *templatesDir == "" || *assetsDir == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	static := StaticFileHandler(path.Join(*assetsDir))
	http.Handle("/", NewUIServer())
	http.Handle("/css/", static)
	http.Handle("/images/", static)
	http.Handle("/js/", static)
	http.HandleFunc("/connect", connect)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
