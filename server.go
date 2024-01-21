package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"socialmedia/contract"
	"socialmedia/service"
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
	svc      *service.Svc
}

func newServer(logger *log.Logger, svc *service.Svc, s *settings.ServerSettings) (*server, error) {
	return &server{
		logger:   logger,
		settings: s,
		svc:      svc,
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
	rq := &contract.RegisterUserRQ{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(rq)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	if rq.Name == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("name isn't set"))
		return
	}

	if rq.Surname == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("surname isn't set"))
		return
	}

	if rq.BirthDate == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("birthdate isn't set"))
		return
	}

	if rq.City == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("city isn't set"))
		return
	}

	if rq.Gender == 0 {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("gender isn't set"))
		return
	}

	if rq.Password == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("password isn't set"))
		return
	}

	id, err := s.svc.RegisterUser(rq)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	s.writeJSON(w, http.StatusOK, &contract.RegisterUserRS{UserID: id})
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
