package contract

type RegisterUserRQ struct {
	Password  string `json:"password"`
	Name      string `json:"first_name"`
	Surname   string `json:"second_name"`
	BirthDate string `json:"birthdate"`
	Gender    int    `json:"gender"`
	Interests string `json:"biography"`
	City      string `json:"city"`
}

type RegisterUserRS struct {
	UserID string `json:"user_id"`
}

type UserRS struct {
	ID        string `json:"id"`
	Name      string `json:"first_name"`
	Surname   string `json:"second_name"`
	BirthDate string `json:"birthdate"`
	Gender    int    `json:"gender"`
	Interests string `json:"biography"`
	City      string `json:"city"`
}

type LoginRQ struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

type LoginRS struct {
	Token string `json:"token"`
}
