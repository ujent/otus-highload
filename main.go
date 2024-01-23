package main

import (
	"log"
	"os"
	"socialmedia/service"
	"socialmedia/settings"
	"socialmedia/storage"

	"github.com/jmoiron/sqlx"
)

func main() {
	logger := log.New(os.Stdout, "otus_socialmedia:", log.LstdFlags|log.Llongfile)

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

	svc, err := service.New(&service.Settings{DB: st, Salt: []byte(s.Salt)})
	if err != nil {
		log.Fatal(err)
	}

	server, err := newServer(logger, svc, s.Server)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
