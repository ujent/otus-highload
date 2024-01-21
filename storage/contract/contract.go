package contract

import "time"

type User struct {
	ID        string
	Password  string
	Name      string
	Surname   string
	BirthDate time.Time
	Gender    int
	Interests string
	City      string
}
