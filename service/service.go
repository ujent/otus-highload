package service

import (
	"socialmedia/contract"
	"socialmedia/hasher"
	"socialmedia/storage"
	dbcontract "socialmedia/storage/contract"
	"time"
)

type Svc struct {
	db   storage.Storage
	salt []byte
}

type Settings struct {
	Salt []byte
	DB   storage.Storage
}

func New(s *Settings) (*Svc, error) {
	return &Svc{db: s.DB, salt: s.Salt}, nil
}

func (svc *Svc) RegisterUser(rq *contract.RegisterUserRQ) (string, error) {
	psw := hasher.GenerateHash([]byte(rq.Password), svc.salt)
	const layout = "2006-01-02"

	birthdate, err := time.Parse(layout, rq.BirthDate)
	if err != nil {
		return "", err
	}

	user := &dbcontract.User{
		Password:  string(psw),
		Name:      rq.Name,
		Surname:   rq.Surname,
		BirthDate: birthdate,
		Gender:    rq.Gender,
		Interests: rq.Interests,
		City:      rq.City,
	}

	id, err := svc.db.AddUser(user)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (svc *Svc) User(id string) (*contract.UserRS, error) {
	user, err := svc.db.User(id)
	if err != nil {
		return nil, err
	}

	res := &contract.UserRS{
		ID:        id,
		Name:      user.Name,
		Surname:   user.Surname,
		BirthDate: user.BirthDate.String(),
		Gender:    user.Gender,
		Interests: user.Interests,
		City:      user.City,
	}

	return res, nil
}
