package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"muon/model"
	"os"
)

var (
	dbFile = flag.String("db", "", "path to database file")
	login = flag.String("login", "", "login for new user")
	email = flag.String("email", "", "email for new user")
)

func main() {
	flag.Parse()

	if *dbFile == "" {
		fmt.Fprintln(os.Stderr, "must provide -db")
		os.Exit(1)
	}

	if db, err := model.Initialize(*dbFile); err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing db:", err)
		os.Exit(1)
	} else {
		defer db.Close()
	}

	if *login == "" || *email == "" {
		fmt.Fprintln(os.Stderr, "must provide -login, -email")
		os.Exit(1)
	}

	pw, err := makepw()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating password:", err)
		os.Exit(1)
	}

	if _, err := model.CreateUser(*login, *email, pw); err != nil {
		fmt.Fprintln(os.Stderr, "Error creating user:", err)
		os.Exit(1)
	}

	fmt.Println("Created user with password:", pw)
}

func makepw() (string, error) {
	// 12 bytes -> 16 bytes base64
	buf := make([]byte, 12)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf), nil
}
