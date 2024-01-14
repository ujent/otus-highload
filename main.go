package main

import (
	"log"
	"os"
	"socialmedia/settings"
	"socialmedia/storage"

	"github.com/jmoiron/sqlx"
)

func main() {
	logger := log.New(os.Stdout, "otus_social:", log.LstdFlags|log.Llongfile)

	s, err := settings.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlx.Connect("mysql", s.DBConnStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	st, err := storage.New(db.DB)
	if err != nil {
		log.Fatal(err)
	}

	server, err := newServer(logger, st, s.Server)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
