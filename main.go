package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

var (
	listenAddr = flag.String("http", ":8080", "")
	dataDir = ""
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... DIR\n", path.Base(os.Args[0]))
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if dataDir = flag.Arg(0); dataDir == "" {
		Usage()
		os.Exit(1)
	}

	static := StaticFileHandler(path.Join(dataDir, "www"))
	http.Handle("/", NewUIServer())
	http.Handle("/css/", static)
	http.Handle("/images/", static)
	http.Handle("/js/", static)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
