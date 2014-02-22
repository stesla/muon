package main

import (
	"flag"
	"log"
	"muon/model"
	"net/http"
	"path"
)

var (
	listenAddr = flag.String("http", ":8080", "")
	dbFile = flag.String("db", "", "database file")
	templatesDir    = flag.String("templates", "./templates", "")
	assetsDir  = flag.String("assets", "./assets", "")
)

func Usage() {
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if *dbFile == "" {
		log.Fatalf("must provide -db")
	}

	if db, err := model.Initialize(*dbFile); err != nil {
		log.Fatalf("error initializing database: %v", err)
	} else {
		defer db.Close()
	}

	if *templatesDir == "" || *assetsDir == "" {
		log.Fatalf("must provide -assets and -templates")
	}

	static := StaticFileHandler(path.Join(*assetsDir))
	http.Handle("/css/", static)
	http.Handle("/images/", static)
	http.Handle("/js/", static)
	http.Handle("/", NewUI(*templatesDir))
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
