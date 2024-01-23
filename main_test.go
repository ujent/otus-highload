package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"socialmedia/contract"
	"socialmedia/service"
	"socialmedia/settings"
	"socialmedia/storage"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
)

const (
	test_salt    = "randomsalt5shf57fhuytglk843dfgdj"
	test_port    = "4000"
	test_connStr = "root:secret@/socialmedia?multiStatements=true&parseTime=true"
)

func TestIntegration(t *testing.T) {
	s := &settings.Settings{
		Server:    &settings.ServerSettings{Port: test_port, Salt: test_salt},
		Salt:      test_salt,
		DBConnStr: test_connStr,
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

	regRQ := &contract.RegisterUserRQ{
		Name:      "Mike",
		Surname:   "Crypton",
		BirthDate: "1990-07-05",
		Gender:    1,
		City:      "Moscow",
		Interests: "everything",
		Password:  "secret",
	}

	id, err := svc.RegisterUser(regRQ)
	if err != nil {
		t.Fatal(err)
	}

	us, err := svc.User(id)
	if err != nil {
		t.Fatal(err)
	}

	isValid, err := svc.IsPswValid(us.ID, regRQ.Password)
	if err != nil {
		t.Fatal(err)
	}

	if !isValid {
		t.Fatalf("psw isn't valid")
	}

	logger := log.New(os.Stdout, "otus_social:", log.LstdFlags|log.Llongfile)

	server, err := newServer(logger, svc, s.Server)
	if err != nil {
		t.Fatal(err)
	}

	token, err := server.generateJWT(id)
	if err != nil {
		t.Fatal(err)
	}

	tokenUserID, err := parseToken(token, s.Salt)
	if err != nil {
		t.Fatal(err)
	}

	if id != tokenUserID {
		t.Fatalf("Wrong userID. Must: %s, has: %s", id, tokenUserID)
	}
}

func parseToken(in, salt string) (userID string, err error) {
	token, err := jwt.Parse(in, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(salt), nil
	})

	if err != nil {
		return "", fmt.Errorf("unauthorized due to error parsing the JWT: %v", err)
	}

	if !token.Valid {
		return "", errors.New("unauthorized due to invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("unauthorized due to invalid claims")
	}

	userID = claims["user"].(string)

	return userID, nil
}
