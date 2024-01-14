package storage

import (
	"database/sql"
	"socialmedia/contract"

	_ "github.com/go-sql-driver/mysql"
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
	ID        string `db:"id"`
	Email     string `db:"email"`
	Password  string `db:"password"`
	Name      string `db:"name"`
	Surname   string `db:"surname"`
	Age       int    `db:"age"`
	Gender    int    `db:"gender"`
	Interests string `db:"interests"`
	City      string `db:"city"`
}

func New(dbPool *sql.DB) (Storage, error) {
	db := sqlx.NewDb(dbPool, "mysql")
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS user
	(id binary(16) DEFAULT (uuid_to_bin(uuid())) NOT NULL PRIMARY KEY, 
		email varchar(255) NOT NULL, 
		password varchar(255) NOT NULL,
		name varchar(255) NOT NULL, 
		surname varchar(255) NOT NULL, 
		age TINYINT NOT NULL, 
		gender TINYINT NOT NULL, 
		city varchar(255) NOT NULL, 
		interests BLOB, 
		UNIQUE (login));`)
	if err != nil {
		return nil, err
	}

	return &storage{db: db}, nil
}

func (st *storage) AddUser(*contract.User) (string, error) {
	return "", nil
}

func (st *storage) User(id string) (*contract.User, error) {
	return nil, nil
}
