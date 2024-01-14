package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"socialmedia/settings"
	"socialmedia/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

var port string

type server struct {
	logger   *log.Logger
	storage  storage.Storage
	settings *settings.ServerSettings
}

func newServer(logger *log.Logger, storage storage.Storage, s *settings.ServerSettings) (*server, error) {
	return &server{
		logger:   logger,
		storage:  storage,
		settings: s,
	}, nil
}

func (s *server) Start() error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	})
	r.Use(cors.Handler)

	r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome"))
		})

		r.Route("/login", func(r chi.Router) {
			r.Post("/", s.login)
		})

		r.Route("/user", func(r chi.Router) {
			r.Post("/register", s.registerUser)
			r.Get("/get/{id}", s.user)
		})

	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		return err
	}

	return nil
}

func (s *server) login(w http.ResponseWriter, r *http.Request) {

}

func (s *server) registerUser(w http.ResponseWriter, r *http.Request) {

}

func (s *server) user(w http.ResponseWriter, r *http.Request) {

}

func (s *server) writeJSON(w http.ResponseWriter, statusCode int, payload interface{}) {

	json, err := json.Marshal(payload)
	if err != nil {
		s.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(json)
}

func (s *server) writeError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Write([]byte(err.Error()))
}
