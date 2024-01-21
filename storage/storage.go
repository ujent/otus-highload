package storage

import (
	"database/sql"
	"socialmedia/storage/contract"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	AddUser(*contract.User) (string, error)
	User(id string) (*contract.User, error)
}

type storage struct {
	db *sqlx.DB
}

type User struct {
	ID        string    `db:"id"`
	Password  string    `db:"password"`
	Name      string    `db:"name"`
	Surname   string    `db:"surname"`
	BirthDate time.Time `db:"birthDate"`
	Gender    int       `db:"gender"`
	Interests string    `db:"interests"`
	City      string    `db:"city"`
}

func New(dbPool *sql.DB) (Storage, error) {
	db := sqlx.NewDb(dbPool, "mysql")
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS user
	(id binary(16) NOT NULL PRIMARY KEY, 
		password varchar(255) NOT NULL,
		name varchar(255) NOT NULL, 
		surname varchar(255) NOT NULL, 
		birthDate DATE NOT NULL, 
		gender TINYINT NOT NULL, 
		city varchar(255) NOT NULL, 
		interests BLOB);`)
	if err != nil {
		return nil, err
	}

	return &storage{db: db}, nil
}

func (st *storage) AddUser(user *contract.User) (string, error) {
	u := uuid.New()
	id, err := u.MarshalBinary()
	if err != nil {
		return "", err
	}

	_, err = st.db.Exec("INSERT INTO user (id, name, surname, password, birthDate, gender, city, interests) VALUES (?,?,?,?,?,?,?,?)", id, user.Name, user.Surname, user.Password, user.BirthDate, user.Gender, user.City, user.Interests)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

func (st *storage) User(id string) (*contract.User, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	uu, err := u.MarshalBinary()
	if err != nil {
		return nil, err
	}

	user := &User{}
	err = st.db.Select(user, "SELECT * FROM user WHERE id = ?", uu)
	if err != nil {
		return nil, err
	}

	res := &contract.User{
		ID:        id,
		Name:      user.Name,
		Surname:   user.Surname,
		Gender:    user.Gender,
		City:      user.City,
		Interests: user.Interests,
		BirthDate: user.BirthDate,
		Password:  user.Password,
	}
	return res, nil
}
