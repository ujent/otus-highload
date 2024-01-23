package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"socialmedia/contract"
	"socialmedia/service"
	"socialmedia/settings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-jwt/jwt"
)

var port string

type server struct {
	logger   *log.Logger
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

			r.Route("/get", func(r chi.Router) {
				r.Use(s.verifyJWT)
				r.Get("/{id}", s.user)
			})
		})

	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		return err
	}

	return nil
}

func (s *server) login(w http.ResponseWriter, r *http.Request) {
	rq := &contract.LoginRQ{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(rq)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}

	if rq.UserID == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("userID isn't set"))
		return
	}

	if rq.Password == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("password isn't set"))
		return
	}

	isValid, err := s.svc.IsPswValid(rq.UserID, rq.Password)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	if !isValid {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("password isn't correct"))
		return
	}

	token, err := s.generateJWT(rq.UserID)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	s.writeJSON(w, http.StatusOK, &contract.LoginRS{Token: token})
}

func (s *server) generateJWT(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["authorized"] = true
	claims["user"] = userID
	tokenString, err := token.SignedString([]byte(s.settings.Salt))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *server) verifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(s.settings.Salt), nil
			})

			if err != nil {
				err := errors.New("unauthorized due to error parsing the JWT")
				s.writeError(w, http.StatusUnauthorized, err)
				return
			}

			if token.Valid {
				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					s.writeError(w, http.StatusUnauthorized, errors.New("unauthorized due to invalid claims"))
					return
				}

				userID := claims["user"].(string)
				ctx := context.WithValue(r.Context(), "userID", userID)

				next.ServeHTTP(w, r.WithContext(ctx))

			} else {
				s.writeError(w, http.StatusUnauthorized, errors.New("unauthorized due to invalid token"))
				return
			}
		} else {
			s.writeError(w, http.StatusUnauthorized, errors.New("unauthorized due to No token in the header"))
			return
		}
	})
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
	id := chi.URLParam(r, "id")
	if id == "" {
		s.writeError(w, http.StatusBadRequest, fmt.Errorf("invalid user id: %s", id))
		return
	}

	res, err := s.svc.User(id)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err)
		return
	}

	s.writeJSON(w, http.StatusOK, res)
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
