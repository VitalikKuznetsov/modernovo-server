package models

type User struct {
	Email       string
	Password    string
	Name        string
	Surname     string
	PhoneNumber string
	DateOfBirth string
}

type UserRegOrLog struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type InfoForUser struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	PhoneNumber string `json:"phonenumber"`
	DateOfBirth string `json:"dateofbirth"`
}
