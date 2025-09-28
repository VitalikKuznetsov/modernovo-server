package models

type User struct {
	Email       string
	Password    string
	Name        string
	PhoneNumber string
}

type UserRegOrLog struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type InfoForUser struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phonenumber"`
}
